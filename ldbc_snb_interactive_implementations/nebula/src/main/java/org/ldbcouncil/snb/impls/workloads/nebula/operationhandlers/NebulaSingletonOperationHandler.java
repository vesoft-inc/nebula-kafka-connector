package org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers;

import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.ClientServerIncompatibleException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import com.vesoft.nebula.client.graph.net.NebulaClient;
import org.ldbcouncil.snb.driver.DbException;
import org.ldbcouncil.snb.driver.Operation;
import org.ldbcouncil.snb.driver.ResultReporter;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.operationhandlers.SingletonOperationHandler;

import java.io.UnsupportedEncodingException;
import java.text.ParseException;
import java.util.Map;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.net.Session;

public abstract class NebulaSingletonOperationHandler<TOperation extends Operation<TOperationResult>, TOperationResult>
        implements SingletonOperationHandler<TOperationResult,TOperation,NebulaDbConnectionState>
{
    public abstract TOperationResult toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException;

    public abstract Map<String, Object> getParameters(NebulaDbConnectionState state, TOperation operation );

    @Override
    public void executeOperation( TOperation operation, NebulaDbConnectionState state,
                                  ResultReporter resultReporter ) throws DbException
    {
        try {
            NebulaClient client = state.getClient();
            String query = getQueryString(state, operation);
            String graphName = state.getGraphName();
            query = query.replace("$graphName", graphName);
            // not implement parameter in session yet
            // final Map<String, Object> parameters = getParameters(state, operation );
            state.logQuery(operation.getClass().getSimpleName(), query);

            final ResultSet resultSet = client.execute(query);

            if (resultSet.rowsSize() > 0) {
                resultReporter.report(1, toResult(resultSet.rowValues(0)), operation);
            } else {
                resultReporter.report(0, null, operation);
            }
        } catch (NoValidSessionException |
                 IOErrorException | ParseException | UnsupportedEncodingException e) {
            throw new RuntimeException(e);
        }
    }
}
