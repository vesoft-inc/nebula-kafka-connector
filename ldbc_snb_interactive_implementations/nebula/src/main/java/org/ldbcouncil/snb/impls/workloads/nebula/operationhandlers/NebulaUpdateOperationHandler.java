package org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.ClientServerIncompatibleException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NotValidConnectionException;
import com.vesoft.nebula.client.graph.net.Session;
import org.ldbcouncil.snb.driver.DbException;
import org.ldbcouncil.snb.driver.Operation;
import org.ldbcouncil.snb.driver.ResultReporter;
import org.ldbcouncil.snb.driver.workloads.interactive.queries.LdbcNoResult;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.operationhandlers.UpdateOperationHandler;

import java.util.Map;

public abstract class NebulaUpdateOperationHandler<TOperation extends Operation<LdbcNoResult>>
        implements UpdateOperationHandler<TOperation,NebulaDbConnectionState>
{
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
            Session session = state.getSession();
            String query = getQueryString(state, operation);
            final Map<String, Object> parameters = getParameters( operation );
            state.logQuery(operation.getClass().getSimpleName(), query);
            final ResultSet resultSet = session.execute(query);
            resultReporter.report( 0, LdbcNoResult.INSTANCE, operation );

        } catch (AuthFailedException | ClientServerIncompatibleException | NotValidConnectionException |
                 IOErrorException e) {
            throw new RuntimeException(e);
        }
    }
}
