package org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import org.ldbcouncil.snb.driver.DbException;
import org.ldbcouncil.snb.driver.Operation;
import org.ldbcouncil.snb.driver.ResultReporter;
import org.ldbcouncil.snb.driver.workloads.interactive.*;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaNewClient;
import org.ldbcouncil.snb.impls.workloads.operationhandlers.UpdateOperationHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

public abstract class NebulaUpdateOperationHandler<TOperation extends Operation<LdbcNoResult>>
        implements UpdateOperationHandler<TOperation,NebulaDbConnectionState>
{
    private static final Logger LOGGER = LoggerFactory.getLogger(NebulaUpdateOperationHandler.class);
    @Override
    public String getQueryString( NebulaDbConnectionState state, TOperation operation )
    {
        return null;
    }

    public abstract Map<String, Object> getParameters(TOperation operation );


    @Override
    public void executeOperation( TOperation operation, NebulaDbConnectionState state,
                                  ResultReporter resultReporter ) throws DbException {

        try {
            NebulaNewClient client = state.getClient(operation.type());
            String query = getQueryString(state, operation);
            String graphName = state.getGraphName();
            query = query.replace("$graphName", graphName);
            // final Map<String, Object> parameters = getParameters( operation );
            state.logQuery(operation.getClass().getSimpleName(), query);
            final ResultSet resultSet = client.getClient().execute(query);
            resultReporter.report( 0, LdbcNoResult.INSTANCE, operation );

        } catch (IOErrorException e) {
            throw new DbException(e);
        }
    }
}
