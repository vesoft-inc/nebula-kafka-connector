package org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers;

import com.vesoft.nebula.client.graph.net.NebulaClient;
import org.ldbcouncil.snb.driver.*;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.nebula.NewNebulaPool;
import org.ldbcouncil.snb.impls.workloads.operationhandlers.MultipleUpdateOperationHandler;
import org.ldbcouncil.snb.driver.workloads.interactive.LdbcNoResult;
import com.vesoft.nebula.client.graph.data.ResultSet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import java.time.LocalTime;
import java.util.List;


public abstract class NebulaMultipleUpdateOperationHandler<TOperation extends Operation<LdbcNoResult>>
        implements MultipleUpdateOperationHandler<TOperation, NebulaDbConnectionState> {

    private static final Logger LOGGER = LoggerFactory.getLogger(NebulaMultipleUpdateOperationHandler.class);

    @Override
    public void executeOperation(TOperation operation, NebulaDbConnectionState state, ResultReporter resultReporter) throws DbException {
        String actualScheduleTime = LocalTime.now().toString();
        NewNebulaPool newPool = null;
        NebulaClient client = null;
        try {
            newPool = state.getPool(operation.type());
            client = newPool.getPool().getClient();
            List<String> queryStrings = getQueryString(state, operation);
            String graphName = state.getGraphName();
            long startTime = System.currentTimeMillis();
            for (String queryString : queryStrings) {
                queryString = queryString.replace("$graphName", graphName);
                state.logQuery(operation.getClass().getSimpleName(), queryString);
                ResultSet resultSet = client.execute(queryString);
                if(state.isEnableQueryInfoLog()){
                    LOGGER.info(String.format("====> sub_update_query=%s, latency=%dus", queryString, resultSet.getLatency()));
                }
                if (!resultSet.isSucceeded()) {
                    LOGGER.error("execute {} failed, {}",
                            operation.getClass().getSimpleName(),
                            resultSet.getErrorMessage());
                }
            }
            long endTime = System.currentTimeMillis();
            if (state.isEnableQueryInfoLog()) {
                LOGGER.info(String.format("====> query type=%s, graphd=%s, execute time=%s, thread_id=%d, response=%dms",
                        operation.getClass().getSimpleName(),
                        newPool.getAddr(),
                        actualScheduleTime,
                        Thread.currentThread().getId(),
                        (endTime - startTime)
                ));
            }
        } catch (Exception e) {
            throw new DbException(e);
        } finally {
            if (newPool != null && client != null) {
                newPool.getPool().returnClient(client);
            }
        }
        resultReporter.report(0, LdbcNoResult.INSTANCE, operation);
    }
}
