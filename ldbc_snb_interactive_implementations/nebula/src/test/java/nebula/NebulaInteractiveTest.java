package nebula;

import org.ldbcouncil.snb.driver.workloads.interactive.db.DummyLdbcSnbInteractiveDb;
import org.ldbcouncil.snb.driver.workloads.interactive.queries.*;
import org.ldbcouncil.snb.driver.workloads.interactive.queries.LdbcShortQuery2PersonPosts;
import org.ldbcouncil.snb.impls.workloads.interactive.InteractiveTest;
import org.ldbcouncil.snb.impls.workloads.nebula.interactive.NebulaInteractiveDb;
import org.junit.Test;

import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class NebulaInteractiveTest extends InteractiveTest {
    public NebulaInteractiveTest()  { super(new NebulaInteractiveDb()); }

    String endpoint = "127.0.0.1:3713";
    String user = "nebula";
    String password = "123";
    String queryDir = "queries";

    @Override
    public Map<String, String> getProperties() {
        Map<String, String> properties = new HashMap<>();
        properties.put("endpoint", endpoint);
        properties.put("user", user);
        properties.put("password", password);
        properties.put("printQueryNames", "true");
        properties.put("printQueryStrings", "true");
        properties.put("printQueryResults", "true");
        properties.put("queryDir", queryDir);
        return properties;
    }

    @Test
    public void testQuery1() throws Exception {
        run(db, new LdbcQuery1(4398046511333L, "johe", LIMIT));
    }

    @Test
    public void testQuery2() throws Exception {
        run(db, new LdbcQuery2(10995116278009L, new Date(1287230400000L), LIMIT));
    }

    @Test
    public void testQuery3() throws Exception {
        long diff = 1277812800000L - 1275393600000L;
        long diffInDays = TimeUnit.DAYS.toDays(diff);
        run(db, new LdbcQuery3a(6597069766734L, "Angola", "Colombia", new Date(1275393600000L), (int) diffInDays, LIMIT));
    }

    @Test
    public void testQuery4() throws Exception {
        long diff = 1277856000000L - 1275350400000L;
        long diffInDays = TimeUnit.DAYS.toDays(diff);
        run(db, new LdbcQuery4(4398046511333L, new Date(1275350400000L), (int) diffInDays, LIMIT));
    }

    @Test
    public void testQuery5() throws Exception {
        run(db, new LdbcQuery5(6597069766734L, new Date(1288612800000L), LIMIT));
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
        run(db, new LdbcQuery9(4398046511268L, new Date(1289908800000L), LIMIT));
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
        run(db, new LdbcQuery13a(8796093022390L, 8796093022357L));
    }

    @Test
    public void testQuery14() throws Exception {
        run(db, new LdbcQuery14a(14L, 27L));
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
    public void testUpdateQuery2() throws Exception {
        run(db, new LdbcInsert2AddPostLike(8796093022239L, 206158430617L, new Date( 1290749436322L ) ) );
    }
}
