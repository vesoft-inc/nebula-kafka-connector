package client

import (
	stderrors "errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
)

var _ = Describe("Client", func() {
	It("Default", func() {
		c := New(WithAddress("127.0.0.1:9669"))
		c1, ok := c.(*defaultClient)
		Expect(ok).To(BeTrue())
		Expect(c1).NotTo(BeNil())
		Expect(c1.addresses).To(Equal([]string{"127.0.0.1:9669"}))
		Expect(c1.user).To(Equal(DefaultUser))
		Expect(c1.password).To(Equal(DefaultPassword))
		Expect(c1.retry).To(Equal(DefaultRetry))
		Expect(c1.concurrencyPerAddress).To(Equal(DefaultConcurrencyPerAddress))
		Expect(c1.queueSize).To(Equal(DefaultQueueSize))
		Expect(c1.queueSize).To(Equal(DefaultQueueSize))
		Expect(c1.chExecuteDataQueue).NotTo(BeNil())
		Expect(c1.done).NotTo(BeNil())
		Expect(c1.logger).NotTo(BeNil())
		Expect(c1.fnNewSession).NotTo(BeNil())
	})

	It("WithAddress", func() {
		c := New(
			WithAddress("127.0.0.1:9669"),
			WithAddress("127.0.0.2:9669,127.0.0.3:9669"),
			WithAddress("127.0.0.4:9669,127.0.0.5:9669", "127.0.0.6:9669"),
			WithUserPassword("newUser", "newPassword"),
			WithReconnectInitialInterval(DefaultReconnectInitialInterval-1),
			WithReconnectInitialInterval(DefaultReconnectInitialInterval+1),
			WithRetry(DefaultRetry-1),
			WithRetry(DefaultRetry+1),
			WithRetryInitialInterval(DefaultRetryInitialInterval-1),
			WithRetryInitialInterval(DefaultRetryInitialInterval+1),
			WithConcurrencyPerAddress(DefaultConcurrencyPerAddress-1),
			WithConcurrencyPerAddress(DefaultConcurrencyPerAddress+1),
			WithQueueSize(DefaultQueueSize-1),
			WithQueueSize(DefaultQueueSize+1),
			WithLogger(logger.NopLogger),
		)
		c1, ok := c.(*defaultClient)
		Expect(ok).To(BeTrue())
		Expect(c1).NotTo(BeNil())
		Expect(c1.addresses).To(Equal([]string{
			"127.0.0.1:9669",
			"127.0.0.2:9669",
			"127.0.0.3:9669",
			"127.0.0.4:9669",
			"127.0.0.5:9669",
			"127.0.0.6:9669",
		}))
		Expect(c1.user).To(Equal("newUser"))
		Expect(c1.password).To(Equal("newPassword"))
		Expect(c1.reconnectInitialInterval).To(Equal(DefaultReconnectInitialInterval + 1))
		Expect(c1.retry).To(Equal(DefaultRetry + 1))
		Expect(c1.retryInitialInterval).To(Equal(DefaultRetryInitialInterval + 1))
		Expect(c1.concurrencyPerAddress).To(Equal(DefaultConcurrencyPerAddress + 1))
		Expect(c1.queueSize).To(Equal(DefaultQueueSize + 1))
		Expect(c1.logger).NotTo(BeNil())
	})

	Describe("Open", func() {
		var (
			ctrl          *gomock.Controller
			mockSession   *MockSession
			mockResultSet *MockResultSet
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockSession = NewMockSession(ctrl)
			mockResultSet = NewMockResultSet(ctrl)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("no addresses", func() {
			c := New()
			err := c.Open()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.ErrNoAddresses))
		})

		It("empty address", func() {
			c := New(WithAddress(""))
			err := c.Open()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.ErrInvalidAddress))
		})

		It("host empty", func() {
			c := New(WithAddress(":9669"))
			err := c.Open()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.ErrInvalidAddress))
		})

		It("port is not a number", func() {
			c := New(WithAddress("127.0.0.1:x"))
			err := c.Open()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.ErrInvalidAddress))
		})

		It("real nebula session", func() {
			c := New(WithAddress("127.0.0.1:0"))
			err := c.Open()
			Expect(err).To(HaveOccurred())
		})

		It("open session failed", func() {
			c := New(
				WithAddress("127.0.0.1:9669"),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)
			mockSession.EXPECT().Open().Return(stderrors.New("test open failed"))
			err := c.Open()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(stderrors.New("test open failed")))
		})

		It("start worker success", func() {
			addresses := []string{"127.0.0.1:9669", "127.0.0.2:9669"}
			c := New(
				WithAddress(addresses...),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)
			mockSession.EXPECT().Close().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			// waiting for all workers to open connection
			time.Sleep(300 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())

			rs, err := c.Execute("test Execute statement")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrClosed))
			Expect(rs).To(BeNil())

			chExecuteResult, ok := c.ExecuteChan("test ExecuteChan statement")
			Expect(ok).To(BeFalse())
			Expect(chExecuteResult).To(BeNil())
		})

		It("start worker failed", func() {
			addresses := []string{"127.0.0.1:9669", "127.0.0.2:9669"}
			c := New(
				WithAddress(addresses...),
				WithReconnectInitialInterval(time.Nanosecond),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)

			var openTimes int64
			var failedOpenTimes = 10
			fnOpen := func() error {
				curr := atomic.AddInt64(&openTimes, 1)
				if curr >= int64(1+DefaultConcurrencyPerAddress)+1 &&
					curr < int64(1+DefaultConcurrencyPerAddress)+1+int64(failedOpenTimes) {
					return stderrors.New("test start worker failed")
				}
				return nil
			}
			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1+DefaultConcurrencyPerAddress)*len(addresses) + failedOpenTimes).DoAndReturn(fnOpen)
			mockSession.EXPECT().Close().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			// waiting for all workers to open connection
			time.Sleep(300 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		It("blocked at ExecuteChan", func() {
			addresses := []string{"127.0.0.1:9669"}
			c := New(
				WithAddress(addresses...),
				WithConcurrencyPerAddress(1),
				WithQueueSize(1),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)

			done := make(chan struct{})
			fnExecute := func(_ string) (ResultSet, error) {
				<-done
				return mockResultSet, nil
			}

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + 1) * len(addresses)).Return(nil)
			mockSession.EXPECT().Execute("test ExecuteChan statement").Times(1).DoAndReturn(fnExecute)
			mockSession.EXPECT().Close().Times((1 + 1) * len(addresses)).Return(nil)

			mockResultSet.EXPECT().IsSucceed().Return(true)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			chExecuteResult, ok := c.ExecuteChan("test ExecuteChan statement")
			Expect(ok).To(BeTrue())
			Expect(chExecuteResult).NotTo(BeNil())

			chExecuteResult, ok = c.ExecuteChan("test ExecuteChan statement")
			Expect(ok).To(BeFalse())
			Expect(chExecuteResult).To(BeNil())

			close(done)

			// waiting for all workers to open connection
			time.Sleep(300 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		It("start worker success", func() {
			addresses := []string{"127.0.0.1:9669", "127.0.0.2:9669"}
			c := New(
				WithAddress(addresses...),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)
			mockSession.EXPECT().Close().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			c1 := c.(*defaultClient)
			close(c1.chExecuteDataQueue)

			// waiting for all workers to open connection
			time.Sleep(300 * time.Millisecond)

			close(c1.done)
			c1.wgSession.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("concurrency", func() {
			var (
				addresses          = []string{"127.0.0.1:9669", "127.0.0.2:9669"}
				executeTimes       = 1000
				executeFailedTimes = 10
				currExecuteTimes1  int64
				currExecuteTimes2  int64
			)

			c := New(
				WithAddress(addresses...),
				WithRetry(executeFailedTimes),
				WithQueueSize(executeTimes*2+executeFailedTimes*2),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)

			fnExecute1 := func(_ string) (ResultSet, error) {
				curr := atomic.AddInt64(&currExecuteTimes1, 1)
				if curr%100 == 0 && curr/100 <= int64(executeFailedTimes) {
					return nil, stderrors.New("execute failed")
				}
				return mockResultSet, nil
			}
			fnExecute2 := func(_ string) (ResultSet, error) {
				curr := atomic.AddInt64(&currExecuteTimes2, 1)
				if curr%100 == 0 && curr/100 <= int64(executeFailedTimes) {
					return nil, stderrors.New("execute failed")
				}
				return mockResultSet, nil
			}

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)
			mockSession.EXPECT().Execute("test Execute statement").Times(executeTimes + executeFailedTimes).DoAndReturn(fnExecute1)
			mockSession.EXPECT().Execute("test ExecuteChan statement").Times(executeTimes + executeFailedTimes).DoAndReturn(fnExecute2)
			mockSession.EXPECT().Close().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)

			mockResultSet.EXPECT().IsSucceed().Times(executeTimes * 2).Return(true)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			var wg sync.WaitGroup
			for i := 0; i < executeTimes; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					rs, err := c.Execute("test Execute statement")
					Expect(err).NotTo(HaveOccurred())
					Expect(rs).NotTo(BeNil())
				}()

				wg.Add(1)
				go func() {
					defer wg.Done()
					chExecuteResult, ok := c.ExecuteChan("test ExecuteChan statement")
					Expect(ok).To(BeTrue())
					executeResult := <-chExecuteResult
					rs, err := executeResult.ResultSet, executeResult.Err
					Expect(err).NotTo(HaveOccurred())
					Expect(rs).NotTo(BeNil())
				}()
			}
			wg.Wait()

			// waiting for all workers to open connection
			time.Sleep(300 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		It("retry case1", func() {
			addresses := []string{"127.0.0.1:9669"}
			c := New(
				WithAddress(addresses...),
				WithRetry(3),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)
			mockSession.EXPECT().Close().AnyTimes().Return(nil)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			// * Case 1: retry no more
			mockSession.EXPECT().Execute("test Execute statement").Times(1).Return(mockResultSet, nil)
			mockResultSet.EXPECT().IsSucceed().Times(1).Return(false)
			mockResultSet.EXPECT().GetError().Times(1).Return(stderrors.New("test error"))
			mockResultSet.EXPECT().IsPermanentError().Times(2).Return(true)

			rs, err := c.Execute("test Execute statement")
			Expect(err).NotTo(HaveOccurred())
			Expect(rs).NotTo(BeNil())
			Expect(rs.IsPermanentError()).To(BeTrue())

			// waiting for all workers to open connection
			time.Sleep(300 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		It("retry case2", func() {
			addresses := []string{"127.0.0.1:9669"}
			c := New(
				WithAddress(addresses...),
				WithRetryInitialInterval(time.Microsecond),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)
			mockSession.EXPECT().Close().AnyTimes().Return(nil)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			retryTimes := DefaultRetry + 10
			var currExecuteTimes int64
			fnIsSucceed := func() bool {
				curr := atomic.AddInt64(&currExecuteTimes, 1)
				return curr > int64(retryTimes)
			}

			// * Case 2. retry as much as possible
			mockSession.EXPECT().Execute("test Execute statement").Times(retryTimes+1).Return(mockResultSet, nil)
			mockResultSet.EXPECT().IsSucceed().Times(retryTimes + 2).DoAndReturn(fnIsSucceed)
			mockResultSet.EXPECT().GetError().Times(retryTimes).Return(stderrors.New("test error"))
			mockResultSet.EXPECT().IsPermanentError().Times(retryTimes).Return(false)
			mockResultSet.EXPECT().IsRetryMoreError().Times(retryTimes).Return(true)

			rs, err := c.Execute("test Execute statement")
			Expect(err).NotTo(HaveOccurred())
			Expect(rs).NotTo(BeNil())
			Expect(rs.IsSucceed()).To(BeTrue())

			// waiting for all workers to open connection
			time.Sleep(300 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		It("retry case3", func() {
			addresses := []string{"127.0.0.1:9669"}
			c := New(
				WithAddress(addresses...),
				WithRetryInitialInterval(time.Microsecond),
				WithNewSessionFunc(func(_ nebula.HostAddress) Session {
					return mockSession
				}),
			)

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)
			mockSession.EXPECT().Close().AnyTimes().Return(nil)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			// * Case 3: retry with limit times
			mockSession.EXPECT().Execute("test Execute statement").Times(DefaultRetry+1).Return(nil, stderrors.New("execute failed"))

			rs, err := c.Execute("test Execute statement")
			Expect(err).To(HaveOccurred())
			Expect(rs).To(BeNil())

			// waiting for all workers to open connection
			time.Sleep(300 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
