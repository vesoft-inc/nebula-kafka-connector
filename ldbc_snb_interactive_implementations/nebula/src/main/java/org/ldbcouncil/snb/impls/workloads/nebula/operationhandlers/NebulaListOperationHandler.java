package org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers;

import com.vesoft.nebula.client.graph.exception.IOErrorException;
import org.ldbcouncil.snb.driver.DbException;
import org.ldbcouncil.snb.driver.Operation;
import org.ldbcouncil.snb.driver.ResultReporter;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.nebula.NebulaNewClient;
import org.ldbcouncil.snb.impls.workloads.operationhandlers.ListOperationHandler;

import java.io.UnsupportedEncodingException;
import java.text.ParseException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import com.vesoft.nebula.client.graph.data.ResultSet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public abstract class NebulaListOperationHandler<TOperation extends Operation<List<TOperationResult>>, TOperationResult>
        implements ListOperationHandler<TOperationResult, TOperation, NebulaDbConnectionState> {

    private static final Logger LOGGER = LoggerFactory.getLogger(NebulaListOperationHandler.class);

    public abstract TOperationResult toResult(ResultSet.Record record) throws ParseException, UnsupportedEncodingException;

    public abstract Map<String, Object> getParameters(NebulaDbConnectionState state, TOperation operation);

    @Override
    public void executeOperation(TOperation operation, NebulaDbConnectionState state,
                                 ResultReporter resultReporter) throws DbException {
        try {
            NebulaNewClient client = state.getClient(operation.type());
            String query = getQueryString(state, operation);
            String graphName = state.getGraphName();
            query = query.replace("$graphName", graphName);
            // not implement parameter in session interface yet.
            // final Map<String, Object> parameters = getParameters(state, operation );
            state.logQuery(operation.getClass().getSimpleName(), query);

            final List<TOperationResult> results = new ArrayList<>();

            final ResultSet resultSet = client.getClient().execute(query);
            while (resultSet.hasNext()) {
                final ResultSet.Record record = resultSet.next();
                results.add(toResult(record));
            }
            resultReporter.report(results.size(), results, operation);
        } catch (IOErrorException | ParseException | UnsupportedEncodingException e) {
            throw new DbException(e);
        }
    }
}
