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
		Expect(c1.fnNewNebulaSession).NotTo(BeNil())
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
			ctrl        *gomock.Controller
			mockSession *MockNebulaSession
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockSession = NewMockNebulaSession(ctrl)
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
				WithNewNebulaSessionFunc(func(_ nebula.HostAddress) NebulaSession {
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
				WithNewNebulaSessionFunc(func(_ nebula.HostAddress) NebulaSession {
					return mockSession
				}),
			)

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)
			mockSession.EXPECT().Close().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)

			err := c.Open()
			Expect(err).NotTo(HaveOccurred())

			// waiting for all workers to open connection
			time.Sleep(100 * time.Millisecond)

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
				WithNewNebulaSessionFunc(func(_ nebula.HostAddress) NebulaSession {
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
			time.Sleep(100 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		It("blocked at ExecuteChan", func() {
			addresses := []string{"127.0.0.1:9669"}
			c := New(
				WithAddress(addresses...),
				WithConcurrencyPerAddress(1),
				WithQueueSize(1),
				WithNewNebulaSessionFunc(func(_ nebula.HostAddress) NebulaSession {
					return mockSession
				}),
			)

			done := make(chan struct{})
			fnExecute := func(_ string) (*nebula.ResultSet, error) {
				<-done
				return &nebula.ResultSet{}, nil
			}

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + 1) * len(addresses)).Return(nil)
			mockSession.EXPECT().Execute("test ExecuteChan statement").Times(1).DoAndReturn(fnExecute)
			mockSession.EXPECT().Close().Times((1 + 1) * len(addresses)).Return(nil)

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
			time.Sleep(100 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		It("start worker success", func() {
			addresses := []string{"127.0.0.1:9669", "127.0.0.2:9669"}
			c := New(
				WithAddress(addresses...),
				WithNewNebulaSessionFunc(func(_ nebula.HostAddress) NebulaSession {
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
			time.Sleep(100 * time.Millisecond)

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
				WithNewNebulaSessionFunc(func(_ nebula.HostAddress) NebulaSession {
					return mockSession
				}),
			)

			fnExecute1 := func(_ string) (*nebula.ResultSet, error) {
				curr := atomic.AddInt64(&currExecuteTimes1, 1)
				if curr%100 == 0 && curr/100 <= int64(executeFailedTimes) {
					return nil, stderrors.New("execute failed")
				}
				return &nebula.ResultSet{}, nil
			}
			fnExecute2 := func(_ string) (*nebula.ResultSet, error) {
				curr := atomic.AddInt64(&currExecuteTimes2, 1)
				if curr%100 == 0 && curr/100 <= int64(executeFailedTimes) {
					return nil, stderrors.New("execute failed")
				}
				return &nebula.ResultSet{}, nil
			}

			// 1 for check and DefaultConcurrencyPerAddress for concurrency per address
			mockSession.EXPECT().Open().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)
			mockSession.EXPECT().Execute("test Execute statement").Times(executeTimes + executeFailedTimes).DoAndReturn(fnExecute1)
			mockSession.EXPECT().Execute("test ExecuteChan statement").Times(executeTimes + executeFailedTimes).DoAndReturn(fnExecute2)
			mockSession.EXPECT().Close().Times((1 + DefaultConcurrencyPerAddress) * len(addresses)).Return(nil)

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
			time.Sleep(100 * time.Millisecond)

			err = c.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
