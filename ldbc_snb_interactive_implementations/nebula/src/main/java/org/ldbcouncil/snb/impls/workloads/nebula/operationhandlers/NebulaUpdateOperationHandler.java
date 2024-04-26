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
            final ResultSet resultSet = client.execute(query);
            if (!resultSet.isSucceeded()) {
                LOGGER.error("execute {} failed, {}",
                        operation.getClass().getSimpleName(),
                        resultSet.getErrorMessage());
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
