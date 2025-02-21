package nebula;

import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.net.NebulaClient;
import org.junit.Assert;
import org.ldbcouncil.snb.driver.workloads.interactive.*;
import org.ldbcouncil.snb.impls.workloads.converter.Converter;
import org.ldbcouncil.snb.impls.workloads.interactive.InteractiveTest;
import org.ldbcouncil.snb.impls.workloads.nebula.interactive.NebulaInteractiveDb;
import org.ldbcouncil.snb.impls.workloads.nebula.converter.NebulaConverter;
import org.junit.Test;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.time.LocalDate;
import java.util.*;
import java.util.concurrent.TimeUnit;

public class NebulaInteractiveTest extends InteractiveTest {
    public NebulaInteractiveTest()  { super(new NebulaInteractiveDb()); }
    protected Converter getConverter() {
        return new NebulaConverter();
    }

    String endpoint = "192.168.8.6:5820";
    String user = "root";
    String password = "Nebula123";
    String queryDir = "queries";
    String requestTimeout = "500";

    @Override
    public Map<String, String> getProperties() {
        Map<String, String> properties = new HashMap<>();
        properties.put("endpoint", endpoint);
        properties.put("user", user);
        properties.put("password", password);
        properties.put("requestTimeout", requestTimeout);
        properties.put("maxSessionSize", "20");
        properties.put("maxSessionWaitTime", "20000");
        properties.put("retryTimes", "3");
        properties.put("intervalTimeBetweenRetrys", "2000");
        properties.put("printQueryNames", "true");
        properties.put("printQueryStrings", "true");
        properties.put("printQueryResults", "true");
        properties.put("graphName", "sf01");
        properties.put("queryDir", queryDir);
        return properties;
    }


    @Test
    public void testConvertTime() throws ParseException {
        System.out.println(NebulaConverter.convertDateToEpoch("2012-09-13"));
    }

    @Test
    public void testLongToDateTime() throws ParseException {
        Long origin = NebulaConverter.convertDateToEpoch(new Date(1331569323280L));
        System.out.println(origin);
        System.out.println(new Date(1331569323280L));

        System.out.println(getConverter().convertDateTime(new Date(1331569323028L)));
        System.out.println(getConverter().convertDateTime(new Date(1347532512905L)));

    }

    @Test
    public void testDateToLong(){

        NebulaConverter converter = new NebulaConverter();
        Date date = new Date(492480000000L);

        try {
            long value = NebulaConverter.convertDateToEpoch(converter.convertDate(date));
            NebulaClient client = NebulaClient.builder("192.168.15.8:9669","root","nebula").build();
            ResultSet    res    = client.execute("use sf10 match(v:Person) where v.id=10995116304675 return v.birthday");
            while(res.hasNext()){
                ResultSet.Record record        = res.next();
                LocalDate        birthday      = record.get(0).asDate();
                String           formatPattern = "yyyy-MM-dd";
                SimpleDateFormat sdf           = new SimpleDateFormat(formatPattern);
                String           formattedDate = sdf.format(converter.convertLocalDateToDate(birthday));
                System.out.println(formattedDate);
            }
            Assert.assertEquals(492480000000L, NebulaConverter.convertDateToEpoch(converter.convertDate(date)));
        }catch (Exception e){
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testQuery1() throws Exception {
        run(db, new LdbcQuery1(4398046511333L, "johe", LIMIT));
    }

    @Test
    public void testQuery2() throws Exception {
        run(db, new LdbcQuery2(10995116278009L, new Date(1333411200000L), LIMIT));
    }

    @Test
    public void testQuery3() throws Exception {
        long diffInDays = 34;
        run(db, new LdbcQuery3(4398046512362L, "Uruguay", "Puerto_Rico", new Date(1325376000000L), (int) diffInDays, LIMIT));
    }

    @Test
    public void testQuery4() throws Exception {
        long diffInDays = 28;
        run(db, new LdbcQuery4(24189255811906L, new Date(1343779200000L), (int) diffInDays, LIMIT));
    }

    @Test
    public void testQuery5() throws Exception {
        run(db, new LdbcQuery5(6597069768287L, new Date(1342483200000L), LIMIT));
    }


    @Test
    public void testQuery6() throws Exception {
        run(db, new LdbcQuery6(4398046511333L, "Carl_Gustaf_Emil_Mannerheim", LIMIT));
    }

    @Test
    public void testQuery7() throws Exception {
        run(db, new LdbcQuery7(4398046511268L, LIMIT));
    }

    @Test
    public void testQuery8() throws Exception {
        run(db, new LdbcQuery8(143L, LIMIT));
    }

    @Test
    public void testQuery9() throws Exception {
        run(db, new LdbcQuery9(454L, new Date(1354060800000L), LIMIT));
    }

    @Test
    public void testQuery10() throws Exception {
        run(db, new LdbcQuery10(4398046511333L, 5, LIMIT));
    }

    @Test
    public void testQuery11() throws Exception {
        run(db, new LdbcQuery11(10995116277918L, "Hungary", 2011, LIMIT));
    }

    @Test
    public void testQuery12() throws Exception {
        run(db, new LdbcQuery12(10995116278009L, "Monarch", LIMIT));
    }

    @Test
    public void testQuery13() throws Exception {
        run(db, new LdbcQuery13(8796093022390L, 8796093022357L));
    }

    @Test
    public void testQuery14() throws Exception {
        run(db, new LdbcQuery14(14L, 27L));
    }

    @Test
    public void testShortQuery1() throws Exception {
        run(db, new LdbcShortQuery1PersonProfile(10995116277794L));
    }

    @Test
    public void testShortQuery2() throws Exception {
        run(db, new LdbcShortQuery2PersonPosts(933L, LIMIT));
    }

    @Test
    public void testShortQuery3() throws Exception {
        run(db, new LdbcShortQuery3PersonFriends(10995116277794L));
    }

    @Test
    public void testShortQuery4() throws Exception {
        run(db, new LdbcShortQuery4MessageContent(206158431836L));
    }

    @Test
    public void testShortQuery5() throws Exception {
        run(db, new LdbcShortQuery5MessageCreator(206158431836L));
    }

    @Test
    public void testShortQuery6() throws Exception {
        run(db, new LdbcShortQuery6MessageForum(206158431836L));
    }

    @Test
    public void testShortQuery7() throws Exception {
        run(db, new LdbcShortQuery7MessageReplies(206158432794L));
    }

    @Test
    public void testUqdateQuery1() throws Exception {
        long personId = 8796093022239L;
        String personFirstName = "wang";
        String personLastName = "jian";
        String gender = "male";
        Date birthDay = new Date(1290749436322L);
        Date creationDate = new Date(1347529090363L);
        String locationIp = "1.12.242.179";
        String browserUsed = "Opera";
        long cityId = 8796093022239L;
        List<String> languages = new ArrayList<>();
        languages.add("chinese");
        languages.add("english");
        List<String> emails = new ArrayList<>();
        emails.add("hello@qq.com");
        emails.add("world@qq.com");
        List<Long> tagIds = new ArrayList<>();
        tagIds.add(123L);
        tagIds.add(234L);
        List<LdbcUpdate1AddPerson.Organization> studyAt = new ArrayList<>();
        LdbcUpdate1AddPerson.Organization university1 = new LdbcUpdate1AddPerson.Organization(2796093022239L, 2001);
        LdbcUpdate1AddPerson.Organization university2 = new LdbcUpdate1AddPerson.Organization(5796093022239L, 2005);
        studyAt.add(university1);
        studyAt.add(university2);
        List<LdbcUpdate1AddPerson.Organization> workAt = new ArrayList<>();
        LdbcUpdate1AddPerson.Organization company1 = new LdbcUpdate1AddPerson.Organization(8796093022239L, 2003);
        LdbcUpdate1AddPerson.Organization company2 = new LdbcUpdate1AddPerson.Organization(1796093022249L, 2004);
        workAt.add(company1);
        workAt.add(company2);
        run(db, new LdbcUpdate1AddPerson(personId, personFirstName, personLastName, gender, birthDay, creationDate, locationIp, browserUsed, cityId, languages, emails, tagIds, studyAt, workAt));
    }

    @Test
    public void testUpdateQuery2() throws Exception {
        run(db, new LdbcUpdate2AddPostLike(8796093022239L, 206158430617L, new Date( 1290749436322L ) ) );
    }

    @Test
    public void testUpdateQuery3() throws Exception {
        run(db, new LdbcUpdate3AddCommentLike(8796093022239L, 206158430617L, new Date( 1290749436322L ) ));
    }

    @Test
    public void testUpdateQuery4() throws Exception {
        List<Long> tagIds = new ArrayList<>();
        run(db, new LdbcUpdate4AddForum(1099511997932L, "tile", new Date(1347529090363L), 206158430617L, tagIds));
    }

    @Test
    public void testUpdateQuery5() throws Exception {
        run(db, new LdbcUpdate5AddForumMembership(8796093022239L, 206158430617L, new Date( 1290749436322L )));
    }

    @Test
    public void testUpdateQuery6() throws Exception {
        List<Long> tagIds = new ArrayList<>();
        tagIds.add(123L);
        tagIds.add(234L);
        run(db, new LdbcUpdate6AddPost(8796093022239L, "image", new Date(1290749436322L), "1.12.242.179", "Opera", "chinese", "roflol", 6, 206158430617L, 8796093022239L, 206158432794L, tagIds));
    }

    @Test
    public void testUpdateQuery7() throws Exception {
        List<Long> tagIds = new ArrayList<>();
        tagIds.add(123L);
        tagIds.add(234L);
        run(db, new LdbcUpdate7AddComment(1099511997932L, new Date(1347529090363L),"1.12.242.179","Opera","roflol", 6, 1L, -1L, 26388279068220L, 1099511997926L, tagIds));
    }

    @Test
    public void testUpdateQuery8() throws Exception {
        run(db, new LdbcUpdate8AddFriendship(8796093022239L, 206158430617L, new Date( 1290749436322L )));
    }
}
