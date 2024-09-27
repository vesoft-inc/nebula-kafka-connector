package org.ldbcouncil.snb.impls.workloads.nebula.converter;

import org.ldbcouncil.snb.driver.workloads.interactive.LdbcQuery1Result;
import org.ldbcouncil.snb.driver.workloads.interactive.LdbcUpdate1AddPerson;
import org.ldbcouncil.snb.impls.workloads.converter.Converter;

import java.text.ParseException;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.List;
import java.util.stream.Collectors;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.TimeZone;


public class NebulaConverter extends Converter {

    final static String DATETIME_FORMAT = "yyyy-MM-dd'T'HH:mm:ss.SSS";
    final static String DATE_FORMAT = "yyyy-MM-dd";

    public String convertOrganisations(List<LdbcUpdate1AddPerson.Organization> values) {
        String res = "[";
        res += values
                .stream()
                .map(v -> "[" + v.getOrganizationId() + ", " + v.getYear() + "]")
                .collect(Collectors.joining(", "));
        res += "]";
        return res;
    }

    public static List<LdbcQuery1Result.Organization> asOrganization(List<List<Object>> value){
        List<LdbcQuery1Result.Organization> orgs = new ArrayList<>();
        for (List<Object> list : value) {
            orgs.add(new LdbcQuery1Result.Organization((String)list.get(0), ((Long) list.get(1)).intValue(), (String)list.get(2)));
        }
        return orgs;
    }

    @Override
    public  String convertStringList(List<String> values) {
        String str = String.join(",", values);
        return "'" + str.replace("'", "\\'") + "'";
    }

    @Override
    public String convertDateTime(Date date) {
        final SimpleDateFormat sdf = new SimpleDateFormat(DATETIME_FORMAT);
        sdf.setTimeZone(TimeZone.getTimeZone("Etc/GMT+0"));
        return "'" + sdf.format(date) + "'";
    }

    @Override
    public String convertDate(Date date) {
        final SimpleDateFormat sdf = new SimpleDateFormat(DATE_FORMAT);
        sdf.setTimeZone(TimeZone.getTimeZone("GMT"));
        return "'" + sdf.format(date) + "'";
    }


    public Date convertLocalDateToDate(LocalDate localDate){
        ZoneId zoneId = ZoneId.of("GMT");
        return Date.from(localDate.atStartOfDay(zoneId).toInstant());
    }
    public static long convertDateTimeToEpoch(LocalDateTime dateTime) {
        ZoneId zoneId = ZoneId.of("GMT");
        return dateTime.atZone(zoneId).toInstant().toEpochMilli();
    }

    public static long convertDateToEpoch(Date date) throws ParseException {
        TimeZone zone = TimeZone.getTimeZone("GMT");
        Calendar calendar = Calendar.getInstance(zone);
        calendar.setTime(date);
         calendar.getTimeInMillis();
         return date.toInstant().toEpochMilli();
    }

    public static long convertDateToEpoch(String date) throws ParseException {
        final SimpleDateFormat sdf = new SimpleDateFormat(DATE_FORMAT);
        sdf.setTimeZone(TimeZone.getTimeZone("GMT"));
        return sdf.parse(date.substring(1, date.length() - 1)).toInstant().toEpochMilli();
    }

    public static int convertStartAndEndDateToLatency(LocalDateTime from, LocalDateTime to) throws ParseException {
        long fromEpoch = convertDateTimeToEpoch(from);
        long toEpoch = convertDateTimeToEpoch(to);
        return (int)((toEpoch - fromEpoch) / 1000 / 60);
    }
}
