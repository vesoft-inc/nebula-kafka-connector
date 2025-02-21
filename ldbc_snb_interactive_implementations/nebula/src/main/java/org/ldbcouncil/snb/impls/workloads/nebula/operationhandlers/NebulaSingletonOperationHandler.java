package org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers;

import com.vesoft.nebula.driver.graph.net.NebulaClient;
import org.ldbcouncil.snb.driver.DbException;
import org.ldbcouncil.snb.driver.Operation;
import org.ldbcouncil.snb.driver.ResultReporter;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.nebula.NewNebulaPool;
import org.ldbcouncil.snb.impls.workloads.nebula.converter.NebulaConverter;
import org.ldbcouncil.snb.impls.workloads.operationhandlers.SingletonOperationHandler;

import java.io.UnsupportedEncodingException;
import java.text.ParseException;
import java.time.LocalTime;
import java.util.Map;

import com.vesoft.nebula.driver.graph.data.ResultSet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public abstract class NebulaSingletonOperationHandler<TOperation extends Operation<TOperationResult>, TOperationResult>
        implements SingletonOperationHandler<TOperationResult, TOperation, NebulaDbConnectionState> {
    private static final Logger LOGGER = LoggerFactory.getLogger(NebulaSingletonOperationHandler.class);

    protected NebulaConverter converter = new NebulaConverter();

    public abstract TOperationResult toResult(ResultSet.Record record) throws ParseException, UnsupportedEncodingException;

    public abstract Map<String, Object> getParameters(NebulaDbConnectionState state, TOperation operation);

    @Override
    public void executeOperation(TOperation operation, NebulaDbConnectionState state,
                                 ResultReporter resultReporter) throws DbException {
        String actualScheduleTime = LocalTime.now().toString();
        NewNebulaPool newPool = null;
        NebulaClient client = null;
        String query = null;
        try {
            newPool = state.getPool(operation.type());
            client = newPool.getPool().getClient();
            query = getQueryString(state, operation);
            String graphName = state.getGraphName();
            query = query.replace("$graphName", graphName);
            // not implement parameter in session yet
            // final Map<String, Object> parameters = getParameters(state, operation );
            state.logQuery(operation.getClass().getSimpleName(), query);

            long startTime = System.currentTimeMillis();
            final ResultSet resultSet = client.execute(query);
            long endTime = System.currentTimeMillis();
            if (!resultSet.isSucceeded()) {
                LOGGER.error("execute {} failed, {}, gql:{}, session id:{}",
                        operation.getClass().getSimpleName(),
                        resultSet.getErrorMessage(),
                             query,
                             client.getSessionId());
            }
            if (state.isEnableQueryInfoLog()) {
                LOGGER.info(String.format("====> query=%s", query));
                LOGGER.info(String.format("====> query type=%s, graphd=%s, execute time=%s, thread_id=%d, latency=%dus, response=%dms",
                        operation.getClass().getSimpleName(),
                        newPool.getAddr(),
                        actualScheduleTime,
                        Thread.currentThread().getId(),
                        resultSet.getLatency(),
                        (endTime - startTime)
                ));
            }
            if (resultSet.hasNext()) {
                resultReporter.report(1, toResult(resultSet.next()), operation);
            } else {
                resultReporter.report(0, null, operation);
            }
        } catch (Exception e) {
            LOGGER.error("====> query {} error", query, e);
            throw new DbException(e);
        } finally {
            if (newPool != null && client != null) {
                newPool.getPool().returnClient(client);
            }
        }
    }
}
