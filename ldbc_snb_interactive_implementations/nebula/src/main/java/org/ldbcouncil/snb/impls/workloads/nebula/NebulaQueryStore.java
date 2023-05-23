package org.ldbcouncil.snb.impls.workloads.nebula;

import java.util.Calendar;
import java.util.Date;
import java.util.Map;
import java.util.List;
import java.util.ArrayList;
import java.util.TimeZone;

import com.google.common.collect.ImmutableMap;

import org.ldbcouncil.snb.driver.DbException;
import org.ldbcouncil.snb.driver.workloads.interactive.*;
import org.ldbcouncil.snb.impls.workloads.QueryStore;
import org.ldbcouncil.snb.impls.workloads.QueryType;
import org.ldbcouncil.snb.impls.workloads.converter.Converter;
import org.ldbcouncil.snb.impls.workloads.nebula.converter.NebulaConverter;


public class NebulaQueryStore extends QueryStore
{
    public NebulaQueryStore( String path )  throws DbException
    {
        super(path, ".gql");
    }

    @Override
    protected Converter getConverter() {
        return new NebulaConverter();
    }

    static protected Date addDays( Date startDate, int days )
    {
        final Calendar cal = Calendar.getInstance(TimeZone.getTimeZone("GMT"));
        cal.setTime( startDate );
        cal.add( Calendar.DATE, days );
        return cal.getTime();
    }

    /**
     * The maps are overriden here to return maps with Java types
     * instead of strings. This speeds up querying the Nebula instance
     * as it reuses the query plans.
     */

    // Complex queries

    @Override
    public Map<String, Object> getQuery1Map(LdbcQuery1 operation) {
        return new ImmutableMap.Builder<String, Object>()
                .put(LdbcQuery1.PERSON_ID, operation.getPersonIdQ1())
                .put(LdbcQuery1.FIRST_NAME, operation.getFirstName())
                .build();
    }

    @Override
    public Map<String, Object> getQuery2Map(LdbcQuery2 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery2.PERSON_ID, operation.getPersonIdQ2() )
        .put(LdbcQuery2.MAX_DATE, getConverter().convertDateTime(operation.getMaxDate()) )
        .put( LdbcQuery2.LIMIT, operation.getLimit() )
        .build();
    }

    @Override
    public Map<String, Object> getQuery3Map(LdbcQuery3 operation) {
        final Date endDate = addDays( operation.getStartDate(), operation.getDurationDays() );
        return new ImmutableMap.Builder<String, Object>()
        .put( LdbcQuery3.PERSON_ID, operation.getPersonIdQ3() )
        .put( LdbcQuery3.COUNTRY_X_NAME, operation.getCountryXName() )
        .put( LdbcQuery3.COUNTRY_Y_NAME, operation.getCountryYName())
        .put( LdbcQuery3.START_DATE, getConverter().convertDateTime(operation.getStartDate()) )
        .put( "endDate", getConverter().convertDateTime(endDate))
        .put( LdbcQuery3.LIMIT, operation.getLimit() )
        .build();
    }

    @Override
    public Map<String, Object> getQuery4Map(LdbcQuery4 operation) {
        final Date endDate = addDays( operation.getStartDate(), operation.getDurationDays() );
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery4.PERSON_ID, operation.getPersonIdQ4())
        .put(LdbcQuery4.START_DATE, getConverter().convertDateTime(operation.getStartDate()))
        .put( "endDate", getConverter().convertDateTime(endDate) )
        .put( LdbcQuery4.LIMIT, operation.getLimit() )
        .build();
    }

    @Override
    public Map<String, Object> getQuery5Map(LdbcQuery5 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery5.PERSON_ID, operation.getPersonIdQ5())
        .put(LdbcQuery5.MIN_DATE, getConverter().convertDateTime(operation.getMinDate()))
        .put( LdbcQuery5.LIMIT, operation.getLimit() )
        .build();
    }

    @Override
    public Map<String, Object> getQuery6Map(LdbcQuery6 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery6.PERSON_ID, operation.getPersonIdQ6())
        .put(LdbcQuery6.TAG_NAME, operation.getTagName())
        .build();
    }

    @Override
    public Map<String, Object> getQuery7Map(LdbcQuery7 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery7.PERSON_ID, operation.getPersonIdQ7())
        .build();
    }

    @Override
    public Map<String, Object> getQuery8Map(LdbcQuery8 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery8.PERSON_ID, operation.getPersonIdQ8())
        .build();
    }

    @Override
    public Map<String, Object> getQuery9Map(LdbcQuery9 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery9.PERSON_ID, operation.getPersonIdQ9())
        .put(LdbcQuery9.MAX_DATE, getConverter().convertDateTime(operation.getMaxDate()) )
        .put( LdbcQuery9.LIMIT, operation.getLimit() )
        .build();
    }

    @Override
    public Map<String, Object> getQuery10Map(LdbcQuery10 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery10.PERSON_ID, operation.getPersonIdQ10())
        .put(LdbcQuery10.MONTH, operation.getMonth())
        .build();
    }

    @Override
    public Map<String, Object> getQuery11Map(LdbcQuery11 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery11.PERSON_ID, operation.getPersonIdQ11())
        .put(LdbcQuery11.COUNTRY_NAME, operation.getCountryName())
        .put(LdbcQuery11.WORK_FROM_YEAR, operation.getWorkFromYear())
        .build();
    }

    @Override
    public Map<String, Object> getQuery12Map(LdbcQuery12 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery12.PERSON_ID, operation.getPersonIdQ12())
        .put(LdbcQuery12.TAG_CLASS_NAME, operation.getTagClassName())
        .build();
    }

    @Override
    public Map<String, Object> getQuery13Map (LdbcQuery13 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery13.PERSON1_ID, operation.getPerson1IdQ13StartNode())
        .put(LdbcQuery13.PERSON2_ID, operation.getPerson2IdQ13EndNode())
        .build();
    }

    @Override
    public Map<String, Object> getQuery14Map (LdbcQuery14 operation) {
        return new ImmutableMap.Builder<String, Object>()
        .put(LdbcQuery14.PERSON1_ID, operation.getPerson1IdQ14StartNode())
        .put(LdbcQuery14.PERSON2_ID, operation.getPerson2IdQ14EndNode())
        .build();
    }

    // Short queries

    @Override
    public Map<String, Object> getShortQuery1PersonProfileMap(LdbcShortQuery1PersonProfile operation) {
        return ImmutableMap.of(LdbcShortQuery1PersonProfile.PERSON_ID, operation.getPersonIdSQ1());
    }

    @Override
    public Map<String, Object> getShortQuery2PersonPostsMap(LdbcShortQuery2PersonPosts operation) {
        return ImmutableMap.of(LdbcShortQuery2PersonPosts.PERSON_ID, operation.getPersonIdSQ2());
    }

    @Override
    public Map<String, Object> getShortQuery3PersonFriendsMap(LdbcShortQuery3PersonFriends operation) {
        return ImmutableMap.of(LdbcShortQuery3PersonFriends.PERSON_ID, operation.getPersonIdSQ3());
    }

    @Override
    public Map<String, Object> getShortQuery4MessageContentMap(LdbcShortQuery4MessageContent operation) {
        return ImmutableMap.of(LdbcShortQuery4MessageContent.MESSAGE_ID, operation.getMessageIdContent());
    }

    @Override
    public Map<String, Object> getShortQuery5MessageCreatorMap(LdbcShortQuery5MessageCreator operation) {
        return ImmutableMap.of(LdbcShortQuery5MessageCreator.MESSAGE_ID, operation.getMessageIdCreator());
    }

    @Override
    public Map<String, Object> getShortQuery6MessageForumMap(LdbcShortQuery6MessageForum operation) {
        return ImmutableMap.of(LdbcShortQuery6MessageForum.MESSAGE_ID, operation.getMessageForumId());
    }

    @Override
    public Map<String, Object> getShortQuery7MessageRepliesMap(LdbcShortQuery7MessageReplies operation) {
        return ImmutableMap.of(LdbcShortQuery7MessageReplies.MESSAGE_ID, operation.getMessageRepliesId());
    }

    @Override
    public String getUpdate2(LdbcUpdate2AddPostLike operation) {
        return prepare(
                QueryType.InteractiveUpdate2,
                ImmutableMap.of(
                        LdbcUpdate2AddPostLike.PERSON_ID, operation.getPersonId(),
                        LdbcUpdate2AddPostLike.POST_ID, operation.getPostId(),
                        LdbcUpdate2AddPostLike.CREATION_DATE, getConverter().convertDateTime(operation.getCreationDate())
                )
        );
    }

    @Override
    public String getUpdate3(LdbcUpdate3AddCommentLike operation) {
        return prepare(
                QueryType.InteractiveUpdate3,
                ImmutableMap.of(
                        LdbcUpdate3AddCommentLike.PERSON_ID, operation.getPersonId(),
                        LdbcUpdate3AddCommentLike.COMMENT_ID, operation.getCommentId(),
                        LdbcUpdate3AddCommentLike.CREATION_DATE, getConverter().convertDateTime(operation.getCreationDate())
                )
        );
    }

    @Override
    public String getUpdate5(LdbcUpdate5AddForumMembership operation) {
        return prepare(
                QueryType.InteractiveUpdate5,
                ImmutableMap.of(
                        LdbcUpdate5AddForumMembership.FORUM_ID, operation.getForumId(),
                        LdbcUpdate5AddForumMembership.PERSON_ID, operation.getPersonId(),
                        LdbcUpdate5AddForumMembership.JOIN_DATE, getConverter().convertDateTime(operation.getJoinDate())
                )
        );
    }

    @Override
    public String getUpdate8(LdbcUpdate8AddFriendship operation) {
        return prepare(
                QueryType.InteractiveUpdate8,
                ImmutableMap.of(
                        LdbcUpdate8AddFriendship.PERSON1_ID, operation.getPerson1Id(),
                        LdbcUpdate8AddFriendship.PERSON2_ID, operation.getPerson2Id(),
                        LdbcUpdate8AddFriendship.CREATION_DATE, getConverter().convertDateTime(operation.getCreationDate())
                )
        );
    }

    @Override
    public List<String> getUpdate1Multiple(LdbcUpdate1AddPerson operation) {
        List<String> list = new ArrayList<>();
        list.add(prepare(
                QueryType.InteractiveUpdate1AddPerson,
                new ImmutableMap.Builder<String, Object>()
                        .put(LdbcUpdate1AddPerson.PERSON_ID, operation.getPersonId())
                        .put(LdbcUpdate1AddPerson.PERSON_FIRST_NAME, getConverter().convertString(operation.getPersonFirstName()))
                        .put(LdbcUpdate1AddPerson.PERSON_LAST_NAME, getConverter().convertString(operation.getPersonLastName()))
                        .put(LdbcUpdate1AddPerson.GENDER, getConverter().convertString(operation.getGender()))
                        .put(LdbcUpdate1AddPerson.BIRTHDAY, getConverter().convertDate(operation.getBirthday()))
                        .put(LdbcUpdate1AddPerson.CREATION_DATE, getConverter().convertDateTime(operation.getCreationDate()))
                        .put(LdbcUpdate1AddPerson.LOCATION_IP, getConverter().convertString(operation.getLocationIp()))
                        .put(LdbcUpdate1AddPerson.BROWSER_USED, getConverter().convertString(operation.getBrowserUsed()))
                        .put(LdbcUpdate1AddPerson.LANGUAGES, getConverter().convertStringList(operation.getLanguages()))
                        .put(LdbcUpdate1AddPerson.EMAILS, getConverter().convertStringList(operation.getEmails()))
                        .build()
        ));

        for (LdbcUpdate1AddPerson.Organization organization : operation.getWorkAt()) {
            list.add(prepare(
                    QueryType.InteractiveUpdate1AddPersonCompanies,
                    ImmutableMap.of(
                            LdbcUpdate1AddPerson.PERSON_ID, operation.getPersonId(),
                            "organizationId", organization.getOrganizationId(),
                            "worksFromYear", organization.getYear()
                    )
            ));
        }

        for (long tagId : operation.getTagIds()) {
            list.add(prepare(
                            QueryType.InteractiveUpdate1AddPersonTags,
                            ImmutableMap.of(
                                    LdbcUpdate1AddPerson.PERSON_ID, operation.getPersonId(),
                                    "tagId", tagId)
                    )
            );
        }
        for (LdbcUpdate1AddPerson.Organization organization : operation.getStudyAt()) {
            list.add(prepare(
                    QueryType.InteractiveUpdate1AddPersonUniversities,
                    ImmutableMap.of(
                            LdbcUpdate1AddPerson.PERSON_ID,operation.getPersonId(),
                            "organizationId", organization.getOrganizationId(),
                            "studiesFromYear", organization.getYear()
                    )
            ));
        }
        // add person->Place(city)
        list.add(prepare(
                QueryType.InteractiveUpdate1AddPersonPlace,
                ImmutableMap.of(
                        LdbcUpdate1AddPerson.PERSON_ID,operation.getPersonId(),
                        LdbcUpdate1AddPerson.CITY_ID, operation.getCityId()
                )
        ));
        return list;
    }


    @Override
    public List<String> getUpdate4Multiple(LdbcUpdate4AddForum operation) {
        List<String> list = new ArrayList<>();
        list.add(prepare(
                QueryType.InteractiveUpdate4AddForum,
                ImmutableMap.of(
                        LdbcUpdate4AddForum.FORUM_ID, operation.getForumId(),
                        LdbcUpdate4AddForum.FORUM_TITLE, getConverter().convertString(operation.getForumTitle()),
                        LdbcUpdate4AddForum.CREATION_DATE, getConverter().convertDateTime(operation.getCreationDate())
                )
        ));

        for (long tagId : operation.getTagIds()) {
            list.add(prepare(
                            QueryType.InteractiveUpdate4AddForumTags,
                            ImmutableMap.of(
                                    LdbcUpdate4AddForum.FORUM_ID, operation.getForumId(),
                                    "tagId",tagId)
                    )
            );
        }
        // add Forum-[hasModerator]->Person
        list.add(prepare(
                QueryType.InteractiveUpdate4AddForumPerson,
                ImmutableMap.of(
                        LdbcUpdate4AddForum.FORUM_ID, operation.getForumId(),
                        LdbcUpdate4AddForum.MODERATOR_PERSON_ID, operation.getModeratorPersonId())
        ));
        return list;
    }

    @Override
    public List<String> getUpdate6Multiple(LdbcUpdate6AddPost operation) {
        List<String> list = new ArrayList<>();
        list.add(prepare(
                        QueryType.InteractiveUpdate6AddPost,
                        new ImmutableMap.Builder<String, Object>()
                                .put(LdbcUpdate6AddPost.POST_ID, operation.getPostId())
                                .put(LdbcUpdate6AddPost.IMAGE_FILE, getConverter().convertString(operation.getImageFile()))
                                .put(LdbcUpdate6AddPost.CREATION_DATE, getConverter().convertDateTime(operation.getCreationDate()))
                                .put(LdbcUpdate6AddPost.LOCATION_IP, getConverter().convertString(operation.getLocationIp()))
                                .put(LdbcUpdate6AddPost.BROWSER_USED, getConverter().convertString(operation.getBrowserUsed()))
                                .put(LdbcUpdate6AddPost.LANGUAGE, getConverter().convertString(operation.getLanguage()))
                                .put(LdbcUpdate6AddPost.CONTENT, getConverter().convertString(operation.getContent()))
                                .put(LdbcUpdate6AddPost.LENGTH, operation.getLength())
                                .build()
                )
        );
        for (long tagId : operation.getTagIds()) {
            list.add(prepare(
                            QueryType.InteractiveUpdate6AddPostTags,
                            ImmutableMap.of(
                                    LdbcUpdate6AddPost.POST_ID, operation.getPostId(),
                                    "tagId", tagId))
            );
        }

        // add Post-[isLocalIn]->Place(Country)
        list.add(prepare(
                QueryType.InteractiveUpdate6AddPostPlace,
                ImmutableMap.of(
                        LdbcUpdate6AddPost.POST_ID, operation.getPostId(),
                        LdbcUpdate6AddPost.COUNTRY_ID, operation.getCountryId())
        ));

        // add Post-[hasCreator]->Person
        list.add(prepare(
                QueryType.InteractiveUpdate6AddPostPerson,
                ImmutableMap.of(
                        LdbcUpdate6AddPost.POST_ID, operation.getPostId(),
                        LdbcUpdate6AddPost.AUTHOR_PERSON_ID, operation.getAuthorPersonId())
        ));

        // add Post-[containerOf]->Forum
        list.add(prepare(
                QueryType.InteractiveUpdate6AddPostForum,
                ImmutableMap.of(
                        LdbcUpdate6AddPost.POST_ID, operation.getPostId(),
                        LdbcUpdate6AddPost.FORUM_ID, operation.getForumId())
        ));
        return list;
    }


    @Override
    public List<String> getUpdate7Multiple(LdbcUpdate7AddComment operation) {
        List<String> list = new ArrayList<>();
        list.add(prepare(
                QueryType.InteractiveUpdate7AddComment,
                new ImmutableMap.Builder<String, Object>()
                        .put(LdbcUpdate7AddComment.COMMENT_ID, operation.getCommentId())
                        .put(LdbcUpdate7AddComment.CREATION_DATE, getConverter().convertDateTime(operation.getCreationDate()))
                        .put(LdbcUpdate7AddComment.LOCATION_IP, getConverter().convertString(operation.getLocationIp()))
                        .put(LdbcUpdate7AddComment.BROWSER_USED, getConverter().convertString(operation.getBrowserUsed()))
                        .put(LdbcUpdate7AddComment.CONTENT, getConverter().convertString(operation.getContent()))
                        .put(LdbcUpdate7AddComment.LENGTH, getConverter().convertInteger(operation.getLength()))
                        .build()
        ));
        for (long tagId : operation.getTagIds()) {
            list.add(prepare(
                            QueryType.InteractiveUpdate7AddCommentTags,
                            ImmutableMap.of(
                                    LdbcUpdate7AddComment.COMMENT_ID, operation.getCommentId(),
                                    "tagId", tagId))
            );
        }

        // add Comment-[Comment_reply_Of]->Comment
        if (operation.getReplyToCommentId() != -1) {
            list.add(prepare(
                    QueryType.InteractiveUpdate7AddCommentComment,
                    ImmutableMap.of(
                            LdbcUpdate7AddComment.COMMENT_ID, operation.getCommentId(),
                            LdbcUpdate7AddComment.REPLY_TO_COMMENT_ID, operation.getReplyToCommentId())
            ));
        }

        // add Comment-[Post_reply_Of]->Post
        if (operation.getReplyToPostId() != -1) {
            list.add(prepare(
                    QueryType.InteractiveUpdate7AddCommentPost,
                    ImmutableMap.of(
                            LdbcUpdate7AddComment.COMMENT_ID, operation.getCommentId(),
                           LdbcUpdate7AddComment.REPLY_TO_POST_ID, operation.getReplyToPostId())
            ));
        }

        // add Comment-[isLocatedIn]->Place(Country)
        list.add(prepare(
                QueryType.InteractiveUpdate7AddCommentPlace,
                ImmutableMap.of(
                        LdbcUpdate7AddComment.COMMENT_ID, operation.getCommentId(),
                        LdbcUpdate7AddComment.COUNTRY_ID, operation.getCountryId())
        ));

        // add Comment-[hasCreator]->Person
        list.add(prepare(
                QueryType.InteractiveUpdate7AddCommentPerson,
                ImmutableMap.of(
                        LdbcUpdate7AddComment.COMMENT_ID, operation.getCommentId(),
                        LdbcUpdate7AddComment.AUTHOR_PERSON_ID, operation.getAuthorPersonId())
        ));

        return list;
    }
}
