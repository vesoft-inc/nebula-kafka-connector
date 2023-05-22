package org.ldbcouncil.snb.impls.workloads.nebula.converter;

import org.ldbcouncil.snb.driver.workloads.interactive.LdbcQuery1Result;
import org.ldbcouncil.snb.driver.workloads.interactive.LdbcUpdate1AddPerson;
import org.ldbcouncil.snb.impls.workloads.converter.Converter;

import java.text.ParseException;
import java.util.ArrayList;
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
    public String convertDateTime(Date date) {
        final SimpleDateFormat sdf = new SimpleDateFormat(DATETIME_FORMAT);
        sdf.setTimeZone(TimeZone.getTimeZone("Etc/GMT+0"));
        return "'" + sdf.format(date) + "'";
    }

    public static long convertDateTimesToEpoch(String date) throws ParseException {
        final SimpleDateFormat sdf = new SimpleDateFormat(DATETIME_FORMAT);
        sdf.setTimeZone(TimeZone.getTimeZone("GMT"));
        return sdf.parse(date.substring(0, date.length() - 3)).toInstant().toEpochMilli();
    }

    public static long convertDateToEpoch(String date) throws ParseException {
        final SimpleDateFormat sdf = new SimpleDateFormat(DATE_FORMAT);
        sdf.setTimeZone(TimeZone.getTimeZone("GMT"));
        return sdf.parse(date).toInstant().toEpochMilli();
    }

    public static int convertStartAndEndDateToLatency(String from, String to) throws ParseException {
        long fromEpoch = convertDateTimesToEpoch(from);
        long toEpoch = convertDateTimesToEpoch(to);
        return (int)((toEpoch - fromEpoch) / 1000 / 60);
    }
}
