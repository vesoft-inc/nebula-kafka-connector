package org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers;

import org.ldbcouncil.snb.driver.*;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaNewClient;
import org.ldbcouncil.snb.impls.workloads.operationhandlers.MultipleUpdateOperationHandler;
import org.ldbcouncil.snb.driver.workloads.interactive.LdbcNoResult;
import com.vesoft.nebula.client.graph.data.ResultSet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import java.util.List;


public abstract class NebulaMultipleUpdateOperationHandler<TOperation extends Operation<LdbcNoResult>>
        implements MultipleUpdateOperationHandler<TOperation, NebulaDbConnectionState> {

    private static final Logger LOGGER = LoggerFactory.getLogger(NebulaMultipleUpdateOperationHandler.class);

    @Override
    public void executeOperation(TOperation operation, NebulaDbConnectionState state, ResultReporter resultReporter) throws DbException {
        try {
            NebulaNewClient client = state.getClient(operation.type());
            List<String> queryStrings = getQueryString(state, operation);
            String graphName = state.getGraphName();
            for (String queryString : queryStrings) {
                queryString = queryString.replace("$graphName", graphName);
                state.logQuery(operation.getClass().getSimpleName(), queryString);
                ResultSet result = client.getClient().execute(queryString);
                if (!result.isSucceeded()) {
                    System.out.println(result.getErrorMessage());
                }
            }
        } catch (Exception e) {
            throw new DbException(e);
        }
        resultReporter.report(0, LdbcNoResult.INSTANCE, operation);
    }
}
