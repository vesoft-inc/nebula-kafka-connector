package org.ldbcouncil.snb.impls.workloads.nebula;

import com.google.common.collect.ImmutableMap;
import com.vesoft.nebula.client.graph.data.ValueWrapper;
import org.ldbcouncil.snb.driver.DbException;
import org.ldbcouncil.snb.driver.control.LoggingService;
import org.ldbcouncil.snb.driver.workloads.interactive.queries.*;
import org.ldbcouncil.snb.impls.workloads.QueryType;
import org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers.NebulaListOperationHandler;
import org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers.NebulaSingletonOperationHandler;
import org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers.NebulaUpdateOperationHandler;
import org.ldbcouncil.snb.impls.workloads.nebula.operationhandlers.NebulaMultipleUpdateOperationHandler;
import org.ldbcouncil.snb.impls.workloads.db.BaseDb;

import java.io.UnsupportedEncodingException;
import java.net.UnknownHostException;
import java.text.ParseException;
import java.util.*;
import java.util.Collections;
import java.util.stream.Collectors;

import com.vesoft.nebula.client.graph.data.ResultSet;

public class NebulaDb extends BaseDb<NebulaQueryStore>
{

    @Override
    protected void onInit( Map<String, String> properties, LoggingService loggingService ) throws DbException
    {
        try {
            dcs = new NebulaDbConnectionState<>(properties, new NebulaQueryStore(properties.get("queryDir")));
        } catch (UnknownHostException e) {
            throw new RuntimeException(e);
        }
    }

    // Interactive complex reads

    public static class InteractiveQuery1 extends NebulaListOperationHandler<LdbcQuery1,LdbcQuery1Result>
    {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery1 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery1);
            return  state.getQueryStore().getQuery1(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery1 operation) {
            return state.getQueryStore().getQuery1Map(operation);
        }

        @Override
        public LdbcQuery1Result toResult(ResultSet.Record record ) throws UnsupportedEncodingException {

            List<String> emails = new ArrayList<>();
            if ( !record.get( 8 ).isNull() ) {
                Collections.addAll(emails, record.get( 8 ).asString().split(","));
            }
            List<String> languages = new ArrayList<>();
            if ( !record.get( 9 ).isNull() ) {
                Collections.addAll(languages, record.get( 9 ).asString().split(","));
            }
            List<LdbcQuery1Result.Organization> universities = new ArrayList<>();
            if ( !record.get( 11 ).isNull() ) {
                List<ValueWrapper> valueList = record.get( 11 ).asList();
                for (ValueWrapper list : valueList) {
                    List<ValueWrapper> res = list.asList();
                    universities.add (new LdbcQuery1Result.Organization(res.get( 0 ).asString(), ( int )res.get( 1 ).asLong(), res.get( 2 ).asString()));
                }
            }
            List<LdbcQuery1Result.Organization> companies = new ArrayList<>();
            if ( !record.get( 12 ).isNull() ) {
                List<ValueWrapper> valueList = record.get(12).asList();
                for (ValueWrapper list : valueList) {
                    List<ValueWrapper> res = list.asList();
                    companies.add(new LdbcQuery1Result.Organization(res.get(0).asString(), ( int )res.get(1).asLong(), res.get(2).asString()));
                }
            }

            long friendId = record.get( 0 ).asLong();
            String friendLastName = record.get( 1 ).asString();
            int distanceFromPerson = (int) record.get( 2 ).asLong();
            long friendBirthday = record.get( 3 ).asLong();
            long friendCreationDate = record.get( 4 ).asLong();
            String friendGender = record.get( 5 ).asString();
            String friendBrowserUsed = record.get( 6 ).asString();
            String friendLocationIp = record.get( 7 ).asString();
            String friendCityName = record.get( 10 ).asString();
            return new LdbcQuery1Result(
                    friendId,
                    friendLastName,
                    distanceFromPerson,
                    friendBirthday,
                    friendCreationDate,
                    friendGender,
                    friendBrowserUsed,
                    friendLocationIp,
                    emails,
                    languages,
                    friendCityName,
                    universities,
                    companies );
        }
    }

    public static class InteractiveQuery2 extends NebulaListOperationHandler<LdbcQuery2,LdbcQuery2Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery2 operation) {
            // (TODO) jmq, we not implement parameters in query yet.
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery2);
            return state.getQueryStore().getQuery2(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery2 operation) {
            return state.getQueryStore().getQuery2Map(operation);
        }

        @Override
        public LdbcQuery2Result toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            long messageId = record.get( 3 ).asLong();
            String messageContent = record.get( 4 ).asString();
            long messageCreationDate = record.get( 5 ).asLong();

            return new LdbcQuery2Result(
                    personId,
                    personFirstName,
                    personLastName,
                    messageId,
                    messageContent,
                    messageCreationDate );
        }
    }

    public static class InteractiveQuery3a extends NebulaListOperationHandler<LdbcQuery3a, LdbcQuery3Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery3a operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery3);
            return state.getQueryStore().getQuery3(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery3a operation) {
            return state.getQueryStore().getQuery3Map(operation);
        }

        @Override
        public LdbcQuery3Result toResult(ResultSet.Record record ) throws UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            int xCount = (int) record.get( 3 ).asLong();
            int yCount = (int) record.get( 4 ).asLong();
            int count = (int) record.get( 5 ).asLong();
            return new LdbcQuery3Result(
                    personId,
                    personFirstName,
                    personLastName,
                    xCount,
                    yCount,
                    count );
        }
    }

    public static class InteractiveQuery3b extends NebulaListOperationHandler<LdbcQuery3b, LdbcQuery3Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery3b operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery3);
            return state.getQueryStore().getQuery3(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery3b operation) {
            return state.getQueryStore().getQuery3Map(operation);
        }

        @Override
        public LdbcQuery3Result toResult(ResultSet.Record record ) throws UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            int xCount = (int) record.get( 3 ).asLong();
            int yCount = (int) record.get( 4 ).asLong();
            int count = (int) record.get( 5 ).asLong();
            return new LdbcQuery3Result(
                    personId,
                    personFirstName,
                    personLastName,
                    xCount,
                    yCount,
                    count );
        }
    }

    public static class InteractiveQuery4 extends NebulaListOperationHandler<LdbcQuery4,LdbcQuery4Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery4 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery4);
            return state.getQueryStore().getQuery4(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery4 operation) {
            return state.getQueryStore().getQuery4Map(operation);
        }

        @Override
        public LdbcQuery4Result toResult(ResultSet.Record record ) throws UnsupportedEncodingException {
            String tagName = record.get( 0 ).asString();
            int postCount = (int) record.get( 1 ).asLong();
            return new LdbcQuery4Result( tagName, postCount );
        }
    }

    public static class InteractiveQuery5 extends NebulaListOperationHandler<LdbcQuery5,LdbcQuery5Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery5 operation) {
            return state.getQueryStore().getQuery5(operation);
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery5);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery5 operation) {
            return state.getQueryStore().getQuery5Map(operation);
        }

        @Override
        public LdbcQuery5Result toResult(ResultSet.Record record ) throws UnsupportedEncodingException {
            String forumTitle = record.get( 0 ).asString();
            int postCount = (int) record.get( 1 ).asLong();
            return new LdbcQuery5Result( forumTitle, postCount );
        }
    }

    public static class InteractiveQuery6 extends NebulaListOperationHandler<LdbcQuery6,LdbcQuery6Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery6 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery6);
            return state.getQueryStore().getQuery6(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery6 operation) {
            return state.getQueryStore().getQuery6Map(operation);
        }

        @Override
        public LdbcQuery6Result toResult(ResultSet.Record record ) throws UnsupportedEncodingException {
            String tagName = record.get( 0 ).asString();
            int postCount = (int) record.get( 1 ).asLong();
            return new LdbcQuery6Result( tagName, postCount );
        }
    }

    public static class InteractiveQuery7 extends NebulaListOperationHandler<LdbcQuery7,LdbcQuery7Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery7 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery7);
            return state.getQueryStore().getQuery7(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery7 operation) {
            return state.getQueryStore().getQuery7Map(operation);
        }

        @Override
        public LdbcQuery7Result toResult(ResultSet.Record record ) throws  UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            long likeCreationDate = record.get( 3 ).asLong();
            long messageId = record.get( 4 ).asLong();
            String messageContent = record.get( 5 ).asString();
            int minutesLatency = (int) record.get( 6 ).asLong();
            boolean isNew = record.get( 7 ).asBoolean();
            return new LdbcQuery7Result(
                    personId,
                    personFirstName,
                    personLastName,
                    likeCreationDate,
                    messageId,
                    messageContent,
                    minutesLatency,
                    isNew );
        }
    }

    public static class InteractiveQuery8 extends NebulaListOperationHandler<LdbcQuery8,LdbcQuery8Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery8 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery8);
            return state.getQueryStore().getQuery8(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery8 operation) {
            return state.getQueryStore().getQuery8Map(operation);
        }

        @Override
        public LdbcQuery8Result toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            long commentCreationDate = record.get( 3 ).asLong();
            long commentId = record.get( 4 ).asLong();
            String commentContent = record.get( 5 ).asString();
            return new LdbcQuery8Result(
                    personId,
                    personFirstName,
                    personLastName,
                    commentCreationDate,
                    commentId,
                    commentContent );
        }
    }

    public static class InteractiveQuery9 extends NebulaListOperationHandler<LdbcQuery9,LdbcQuery9Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery9 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery9);
            return state.getQueryStore().getQuery9(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery9 operation) {
            return state.getQueryStore().getQuery9Map(operation);
        }

        @Override
        public LdbcQuery9Result toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            long messageId = record.get( 3 ).asLong();
            String messageContent = record.get( 4 ).asString();
            long messageCreationDate = record.get( 5 ).asLong();
            return new LdbcQuery9Result(
                    personId,
                    personFirstName,
                    personLastName,
                    messageId,
                    messageContent,
                    messageCreationDate );
        }
    }

    public static class InteractiveQuery10 extends NebulaListOperationHandler<LdbcQuery10,LdbcQuery10Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery10 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery10);
            return state.getQueryStore().getQuery10(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery10 operation) {
            return state.getQueryStore().getQuery10Map(operation);
        }

        @Override
        public LdbcQuery10Result toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            int commonInterestScore = (int) record.get( 3 ).asLong();
            String personGender = record.get( 4 ).asString();
            String personCityName = record.get( 5 ).asString();
            return new LdbcQuery10Result(
                    personId,
                    personFirstName,
                    personLastName,
                    commonInterestScore,
                    personGender,
                    personCityName );
        }
    }

    public static class InteractiveQuery11 extends NebulaListOperationHandler<LdbcQuery11,LdbcQuery11Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery11 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery11);
            return state.getQueryStore().getQuery11(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery11 operation) {
            return state.getQueryStore().getQuery11Map(operation);
        }

        @Override
        public LdbcQuery11Result toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            String organizationName = record.get( 3 ).asString();
            int organizationWorkFromYear = (int) record.get( 4 ).asLong();
            return new LdbcQuery11Result(
                    personId,
                    personFirstName,
                    personLastName,
                    organizationName,
                    organizationWorkFromYear );
        }
    }

    public static class InteractiveQuery12 extends NebulaListOperationHandler<LdbcQuery12,LdbcQuery12Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery12 operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery12);
            return state.getQueryStore().getQuery12(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery12 operation) {
            return state.getQueryStore().getQuery12Map(operation);
        }

        @Override
        public LdbcQuery12Result toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String personFirstName = record.get( 1 ).asString();
            String personLastName = record.get( 2 ).asString();
            List<String> tagNames = new ArrayList<>();
            if (!record.get(3).isNull()) {
                List<ValueWrapper> valueList = record.get(3).asList();
                for (ValueWrapper val : valueList) {
                    tagNames.add(val.asString());
                }
            }
            int replyCount = (int) record.get( 4 ).asLong();
            return new LdbcQuery12Result(
                    personId,
                    personFirstName,
                    personLastName,
                    tagNames,
                    replyCount );
        }
    }

    public static class InteractiveQuery13a extends NebulaSingletonOperationHandler<LdbcQuery13a, LdbcQuery13Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery13a operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery13);
            return state.getQueryStore().getQuery13(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery13a operation) {
            return state.getQueryStore().getQuery13Map(operation);
        }

        @Override
        public LdbcQuery13Result toResult(ResultSet.Record record )
        {
            if(record != null) {
                return new LdbcQuery13Result( (int) record.get( 0 ).asLong() );
            } else {
                return new LdbcQuery13Result(-1);
            }
        }
    }

    public static class InteractiveQuery13b extends NebulaSingletonOperationHandler<LdbcQuery13b,LdbcQuery13Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery13b operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery13);
            return state.getQueryStore().getQuery13(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery13b operation) {
            return state.getQueryStore().getQuery13Map(operation);
        }

        @Override
        public LdbcQuery13Result toResult(ResultSet.Record record )
        {
            if(record != null) {
                return new LdbcQuery13Result( (int) record.get( 0 ).asLong() );
            } else {
                return new LdbcQuery13Result(-1);
            }
        }
    }

    public static class InteractiveQuery14a extends NebulaListOperationHandler<LdbcQuery14a, LdbcQuery14Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery14a operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery14);
            return state.getQueryStore().getQuery14(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery14a operation) {
            return state.getQueryStore().getQuery14Map(operation);
        }

        @Override
        public LdbcQuery14Result toResult(ResultSet.Record record ) throws ParseException
        {
            List<Long> personIdsInPath = new ArrayList<>();
            if (!record.get(0).isNull()) {
                List<ValueWrapper> values = record.get(0).asList();
                for (ValueWrapper val : values) {
                    personIdsInPath.add(val.asLong());
                }
            }
            long pathWeight = record.get( 1 ).asLong();
            return new LdbcQuery14Result(
                    personIdsInPath,
                    pathWeight );
        }
    }

    public static class InteractiveQuery14b extends NebulaListOperationHandler<LdbcQuery14b, LdbcQuery14Result>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcQuery14b operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveComplexQuery14);
            return state.getQueryStore().getQuery14(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcQuery14b operation) {
            return state.getQueryStore().getQuery14Map(operation);
        }

        @Override
        public LdbcQuery14Result toResult(ResultSet.Record record ) throws ParseException
        {
            List<Long> personIdsInPath = new ArrayList<>();
            if (!record.get(0).isNull()) {
                List<ValueWrapper> values = record.get(0).asList();
                for (ValueWrapper val : values) {
                    personIdsInPath.add(val.asLong());
                }
            }
            long pathWeight = record.get( 1 ).asLong();
            return new LdbcQuery14Result(
                    personIdsInPath,
                    pathWeight );
        }
    }

    // Interactive short reads
    public static class ShortQuery1PersonProfile extends NebulaSingletonOperationHandler<LdbcShortQuery1PersonProfile,LdbcShortQuery1PersonProfileResult>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcShortQuery1PersonProfile operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveShortQuery1);
            return state.getQueryStore().getShortQuery1PersonProfile(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcShortQuery1PersonProfile operation) {
            return state.getQueryStore().getShortQuery1PersonProfileMap(operation);
        }

        @Override
        public LdbcShortQuery1PersonProfileResult toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            if (record != null){
                String firstName = record.get( 0 ).asString();
                String lastName = record.get( 1 ).asString();
                long birthday = record.get( 2 ).asLong();
                String locationIP = record.get( 3 ).asString();
                String browserUsed = record.get( 4 ).asString();
                long cityId = record.get( 5 ).asLong();
                String gender = record.get( 6 ).asString();
                long creationDate = record.get( 7 ).asLong();
                return new LdbcShortQuery1PersonProfileResult(
                        firstName,
                        lastName,
                        birthday,
                        locationIP,
                        browserUsed,
                        cityId,
                        gender,
                        creationDate );
            }
            else
            {
                return null;
            }

        }
    }

    public static class ShortQuery2PersonPosts extends NebulaListOperationHandler<LdbcShortQuery2PersonPosts,LdbcShortQuery2PersonPostsResult>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcShortQuery2PersonPosts operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveShortQuery2);
            return state.getQueryStore().getShortQuery2PersonPosts(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcShortQuery2PersonPosts operation) {
            return state.getQueryStore().getShortQuery2PersonPostsMap(operation);
        }

        @Override
        public LdbcShortQuery2PersonPostsResult toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            if (record != null){
                long messageId = record.get( 0 ).asLong();
                String messageContent = record.get( 1 ).asString();
                long messageCreationDate = record.get( 2 ).asLong();
                long originalPostId = record.get( 3 ).asLong();
                long originalPostAuthorId = record.get( 4 ).asLong();
                String originalPostAuthorFirstName = record.get( 5 ).asString();
                String originalPostAuthorLastName = record.get( 6 ).asString();
                return new LdbcShortQuery2PersonPostsResult(
                        messageId,
                        messageContent,
                        messageCreationDate,
                        originalPostId,
                        originalPostAuthorId,
                        originalPostAuthorFirstName,
                        originalPostAuthorLastName );
            }
            else
            {
                return null;
            }
        }
    }

    public static class ShortQuery3PersonFriends extends NebulaListOperationHandler<LdbcShortQuery3PersonFriends,LdbcShortQuery3PersonFriendsResult>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcShortQuery3PersonFriends operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveShortQuery3);
            return state.getQueryStore().getShortQuery3PersonFriends(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcShortQuery3PersonFriends operation) {
            return state.getQueryStore().getShortQuery3PersonFriendsMap(operation);
        }

        @Override
        public LdbcShortQuery3PersonFriendsResult toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            long personId = record.get( 0 ).asLong();
            String firstName = record.get( 1 ).asString();
            String lastName = record.get( 2 ).asString();
            long friendshipCreationDate = record.get( 3 ).asLong();
            return new LdbcShortQuery3PersonFriendsResult(
                    personId,
                    firstName,
                    lastName,
                    friendshipCreationDate );
        }
    }

    public static class ShortQuery4MessageContent extends NebulaSingletonOperationHandler<LdbcShortQuery4MessageContent,LdbcShortQuery4MessageContentResult>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcShortQuery4MessageContent operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveShortQuery4);
            return state.getQueryStore().getShortQuery4MessageContent(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcShortQuery4MessageContent operation) {
            return state.getQueryStore().getShortQuery4MessageContentMap(operation);
        }

        @Override
        public LdbcShortQuery4MessageContentResult toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            if (record != null){
            // Pay attention, the spec's and the implementation's parameter orders are different.
            long messageCreationDate = record.get( 0 ).asLong();
            String messageContent = record.get( 1 ).asString();
            return new LdbcShortQuery4MessageContentResult(
                    messageContent,
                    messageCreationDate );
            }
            else{
                return null;
            }

        }
    }

    public static class ShortQuery5MessageCreator extends NebulaSingletonOperationHandler<LdbcShortQuery5MessageCreator,LdbcShortQuery5MessageCreatorResult>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcShortQuery5MessageCreator operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveShortQuery5);
            return state.getQueryStore().getShortQuery5MessageCreator(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcShortQuery5MessageCreator operation) {
            return state.getQueryStore().getShortQuery5MessageCreatorMap(operation);
        }

        @Override
        public LdbcShortQuery5MessageCreatorResult toResult(ResultSet.Record record ) throws UnsupportedEncodingException {
            if (record != null){
                long personId = record.get( 0 ).asLong();
                String firstName = record.get( 1 ).asString();
                String lastName = record.get( 2 ).asString();
                return new LdbcShortQuery5MessageCreatorResult(
                        personId,
                        firstName,
                        lastName );
            }
            else
            {
                return null;
            }

        }
    }

    public static class ShortQuery6MessageForum extends NebulaSingletonOperationHandler<LdbcShortQuery6MessageForum,LdbcShortQuery6MessageForumResult>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcShortQuery6MessageForum operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveShortQuery6);
            return state.getQueryStore().getShortQuery6MessageForum(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcShortQuery6MessageForum operation) {
            return state.getQueryStore().getShortQuery6MessageForumMap(operation);
        }

        @Override
        public LdbcShortQuery6MessageForumResult toResult(ResultSet.Record record ) throws UnsupportedEncodingException {
            if (record != null)
            {
                long forumId = record.get( 0 ).asLong();
                String forumTitle = record.get( 1 ).asString();
                long moderatorId = record.get( 2 ).asLong();
                String moderatorFirstName = record.get( 3 ).asString();
                String moderatorLastName = record.get( 4 ).asString();
                return new LdbcShortQuery6MessageForumResult(
                        forumId,
                        forumTitle,
                        moderatorId,
                        moderatorFirstName,
                        moderatorLastName );
            }
            else
            {
                return null;
            }
        }
    }

    public static class ShortQuery7MessageReplies extends NebulaListOperationHandler<LdbcShortQuery7MessageReplies,LdbcShortQuery7MessageRepliesResult>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcShortQuery7MessageReplies operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveShortQuery7);
            return state.getQueryStore().getShortQuery7MessageReplies(operation);
        }

        @Override
        public Map<String, Object> getParameters(NebulaDbConnectionState state, LdbcShortQuery7MessageReplies operation) {
            return state.getQueryStore().getShortQuery7MessageRepliesMap(operation);
        }

        @Override
        public LdbcShortQuery7MessageRepliesResult toResult(ResultSet.Record record ) throws ParseException, UnsupportedEncodingException {
            long commentId = record.get( 0 ).asLong();
            String commentContent = record.get( 1 ).asString();
            long commentCreationDate = record.get( 2 ).asLong();
            long replyAuthorId = record.get( 3 ).asLong();
            String replyAuthorFirstName = record.get( 4 ).asString();
            String replyAuthorLastName = record.get( 5 ).asString();
            boolean replyAuthorKnowsOriginalMessageAuthor = record.get( 6 ).asBoolean();
            return new LdbcShortQuery7MessageRepliesResult(
                    commentId,
                    commentContent,
                    commentCreationDate,
                    replyAuthorId,
                    replyAuthorFirstName,
                    replyAuthorLastName,
                    replyAuthorKnowsOriginalMessageAuthor );
        }
    }

    // Interactive inserts

    public static class Insert1AddPerson extends NebulaMultipleUpdateOperationHandler<LdbcInsert1AddPerson>
    {
        @Override
        public List<String> getQueryString(NebulaDbConnectionState state, LdbcInsert1AddPerson operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveInsert1);
            // (TODO) jmq, runtime insert is not implemete yet.
            return state.getQueryStore().getInsert1Multiple(operation);
        }
    }

    public static class Insert2AddPostLike extends NebulaUpdateOperationHandler<LdbcInsert2AddPostLike>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcInsert2AddPostLike operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveInsert2);
            return state.getQueryStore().getInsert2(operation);
        }

        @Override
        public Map<String, Object> getParameters( LdbcInsert2AddPostLike operation )
        {
            return ImmutableMap.<String, Object>builder()
                               .put( LdbcInsert2AddPostLike.PERSON_ID, operation.getPersonId() )
                               .put( LdbcInsert2AddPostLike.POST_ID, operation.getPostId() )
                               .put( LdbcInsert2AddPostLike.CREATION_DATE, operation.getCreationDate().getTime() )
                               .build();
        }
    }

    public static class Insert3AddCommentLike extends NebulaUpdateOperationHandler<LdbcInsert3AddCommentLike>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcInsert3AddCommentLike operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveInsert3);
            return state.getQueryStore().getInsert3(operation);
        }

        @Override
        public Map<String, Object> getParameters( LdbcInsert3AddCommentLike operation )
        {
            return ImmutableMap.<String, Object>builder()
                               .put( LdbcInsert3AddCommentLike.PERSON_ID, operation.getPersonId() )
                               .put( LdbcInsert3AddCommentLike.COMMENT_ID, operation.getCommentId() )
                               .put( LdbcInsert3AddCommentLike.CREATION_DATE, operation.getCreationDate().getTime() )
                               .build();
        }
    }

    public static class Insert4AddForum extends NebulaMultipleUpdateOperationHandler<LdbcInsert4AddForum>
    {
        @Override
        public List<String> getQueryString(NebulaDbConnectionState state, LdbcInsert4AddForum operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveInsert4);
            return state.getQueryStore().getInsert4Multiple(operation);
        }
    }

    public static class Insert5AddForumMembership extends NebulaUpdateOperationHandler<LdbcInsert5AddForumMembership>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcInsert5AddForumMembership operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveInsert5);
            return state.getQueryStore().getInsert5(operation);
        }

        @Override
        public Map<String, Object> getParameters( LdbcInsert5AddForumMembership operation )
        {
            return ImmutableMap.<String, Object>builder()
                               .put( LdbcInsert5AddForumMembership.FORUM_ID, operation.getForumId() )
                               .put( LdbcInsert5AddForumMembership.PERSON_ID, operation.getPersonId() )
                               .put( LdbcInsert5AddForumMembership.CREATION_DATE, operation.getCreationDate().getTime() )
                               .build();
        }
    }

    public static class Insert6AddPost extends NebulaMultipleUpdateOperationHandler<LdbcInsert6AddPost>
    {
        @Override
        public List<String> getQueryString(NebulaDbConnectionState state, LdbcInsert6AddPost operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveInsert6);
            return state.getQueryStore().getInsert6Multiple(operation);
        }
    }

    public static class Insert7AddComment extends NebulaMultipleUpdateOperationHandler<LdbcInsert7AddComment>
    {
        @Override
        public List<String> getQueryString(NebulaDbConnectionState state, LdbcInsert7AddComment operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveInsert7);
            return state.getQueryStore().getInsert7Multiple(operation);
        }
    }

    public static class Insert8AddFriendship extends NebulaUpdateOperationHandler<LdbcInsert8AddFriendship>
    {
        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcInsert8AddFriendship operation) {
            // return state.getQueryStore().getParameterizedQuery(QueryType.InteractiveInsert8);
            return state.getQueryStore().getInsert8(operation);
        }

        @Override
        public Map<String, Object> getParameters( LdbcInsert8AddFriendship operation )
        {
            return ImmutableMap.<String, Object>builder()
                               .put( LdbcInsert8AddFriendship.PERSON1_ID, operation.getPerson1Id() )
                               .put( LdbcInsert8AddFriendship.PERSON2_ID, operation.getPerson2Id() )
                               .put( LdbcInsert8AddFriendship.CREATION_DATE, operation.getCreationDate().getTime() )
                               .build();
        }
    }

    // Deletions
    public static class Delete1RemovePerson extends NebulaUpdateOperationHandler<LdbcDelete1RemovePerson> {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcDelete1RemovePerson operation) {
            return state.getQueryStore().getDelete1(operation);
        }

        @Override
        public Map<String, Object> getParameters(LdbcDelete1RemovePerson operation) {
            return ImmutableMap.<String, Object>builder()
                    .put( LdbcDelete1RemovePerson.PERSON_ID, operation.getremovePersonIdD1() )
                    .build();
        }
    }

    public static class Delete2RemovePostLike extends NebulaUpdateOperationHandler<LdbcDelete2RemovePostLike> {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcDelete2RemovePostLike operation) {
            return state.getQueryStore().getDelete2(operation);
        }

        @Override
        public Map<String, Object> getParameters(LdbcDelete2RemovePostLike operation) {
            return ImmutableMap.<String, Object>builder()
                    .put( LdbcDelete2RemovePostLike.PERSON_ID, operation.getremovePersonIdD2() )
                    .put( LdbcDelete2RemovePostLike.POST_ID, operation.getremovePostIdD2() )
                    .build();
        }
    }

    public static class Delete3RemoveCommentLike extends NebulaUpdateOperationHandler<LdbcDelete3RemoveCommentLike> {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcDelete3RemoveCommentLike operation) {
            return state.getQueryStore().getDelete3(operation);
        }

        @Override
        public Map<String, Object> getParameters(LdbcDelete3RemoveCommentLike operation) {
            return ImmutableMap.<String, Object>builder()
                    .put( LdbcDelete3RemoveCommentLike.PERSON_ID, operation.getremovePersonIdD3() )
                    .put( LdbcDelete3RemoveCommentLike.COMMENT_ID, operation.getremoveCommentIdD3() )
                    .build();
        }
    }

    public static class Delete4RemoveForum extends NebulaUpdateOperationHandler<LdbcDelete4RemoveForum> {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcDelete4RemoveForum operation) {
            return state.getQueryStore().getDelete4(operation);
        }

        @Override
        public Map<String, Object> getParameters(LdbcDelete4RemoveForum operation) {
            return ImmutableMap.<String, Object>builder()
                    .put( LdbcDelete4RemoveForum.FORUM_ID, operation.getremoveForumIdD4() )
                    .build();
        }
    }

    public static class Delete5RemoveForumMembership extends NebulaUpdateOperationHandler<LdbcDelete5RemoveForumMembership> {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcDelete5RemoveForumMembership operation) {
            return state.getQueryStore().getDelete5(operation);
        }

        @Override
        public Map<String, Object> getParameters(LdbcDelete5RemoveForumMembership operation) {
            return ImmutableMap.<String, Object>builder()
                    .put( LdbcDelete5RemoveForumMembership.PERSON_ID, operation.getremovePersonIdD5() )
                    .put( LdbcDelete5RemoveForumMembership.FORUM_ID, operation.getremoveForumIdD5() )
                    .build();
        }
    }

    public static class Delete6RemovePostThread extends NebulaUpdateOperationHandler<LdbcDelete6RemovePostThread> {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcDelete6RemovePostThread operation) {
            return state.getQueryStore().getDelete6(operation);
        }

        @Override
        public Map<String, Object> getParameters(LdbcDelete6RemovePostThread operation) {
            return ImmutableMap.<String, Object>builder()
                    .put( LdbcDelete6RemovePostThread.POST_ID, operation.getremovePostIdD6() )
                    .build();
        }
    }

    public static class Delete7RemoveCommentSubthread extends NebulaUpdateOperationHandler<LdbcDelete7RemoveCommentSubthread> {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcDelete7RemoveCommentSubthread operation) {
            return state.getQueryStore().getDelete7(operation);
        }

        @Override
        public Map<String, Object> getParameters(LdbcDelete7RemoveCommentSubthread operation) {
            return ImmutableMap.<String, Object>builder()
                    .put( LdbcDelete7RemoveCommentSubthread.COMMENT_ID, operation.getremoveCommentIdD7() )
                    .build();
        }
    }

    public static class Delete8RemoveFriendship extends NebulaUpdateOperationHandler<LdbcDelete8RemoveFriendship> {

        @Override
        public String getQueryString(NebulaDbConnectionState state, LdbcDelete8RemoveFriendship operation) {
            return state.getQueryStore().getDelete8(operation);
        }

        @Override
        public Map<String, Object> getParameters(LdbcDelete8RemoveFriendship operation) {
            return ImmutableMap.<String, Object>builder()
                    .put( LdbcDelete8RemoveFriendship.PERSON1_ID, operation.getremovePerson1Id() )
                    .put( LdbcDelete8RemoveFriendship.PERSON2_ID, operation.getremovePerson2Id() )
                    .build();
        }
    }

}
