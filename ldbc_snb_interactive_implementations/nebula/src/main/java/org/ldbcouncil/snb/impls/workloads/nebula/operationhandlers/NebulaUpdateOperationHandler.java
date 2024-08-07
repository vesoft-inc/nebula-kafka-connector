package org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.net.NebulaClient;
import org.ldbcouncil.snb.driver.DbException;
import org.ldbcouncil.snb.driver.Operation;
import org.ldbcouncil.snb.driver.ResultReporter;
import org.ldbcouncil.snb.driver.workloads.interactive.*;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.nebula.NewNebulaPool;
import org.ldbcouncil.snb.impls.workloads.operationhandlers.UpdateOperationHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.LocalTime;
import java.util.Map;

public abstract class NebulaUpdateOperationHandler<TOperation extends Operation<LdbcNoResult>>
        implements UpdateOperationHandler<TOperation, NebulaDbConnectionState> {
    private static final Logger LOGGER = LoggerFactory.getLogger(NebulaUpdateOperationHandler.class);

    @Override
    public String getQueryString(NebulaDbConnectionState state, TOperation operation) {
        return null;
    }

    public abstract Map<String, Object> getParameters(TOperation operation);


    @Override
    public void executeOperation(TOperation operation, NebulaDbConnectionState state,
                                 ResultReporter resultReporter) throws DbException {
        String actualScheduleTime = LocalTime.now().toString();
        NewNebulaPool newPool = null;
        NebulaClient client = null;
        try {
            newPool = state.getPool(operation.type());
            client = newPool.getPool().getClient();
            String query = getQueryString(state, operation);
            String graphName = state.getGraphName();
            query = query.replace("$graphName", graphName);
            // final Map<String, Object> parameters = getParameters( operation );
            state.logQuery(operation.getClass().getSimpleName(), query);

            long startTime = System.currentTimeMillis();
            final ResultSet resultSet = client.execute(query);
            long endTime = System.currentTimeMillis();
            if (!resultSet.isSucceeded()) {
                LOGGER.error("execute {} failed, {}, session id:{}",
                        operation.getClass().getSimpleName(),
                        resultSet.getErrorMessage(),
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
            resultReporter.report(0, LdbcNoResult.INSTANCE, operation);
        } catch (Exception e) {
            throw new DbException(e);
        } finally {
            if (newPool != null && client != null) {
                newPool.getPool().returnClient(client);
            }
        }
    }
}
