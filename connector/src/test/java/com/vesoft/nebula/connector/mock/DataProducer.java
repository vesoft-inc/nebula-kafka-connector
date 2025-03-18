/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.mock;

import com.alibaba.fastjson.JSONObject;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.common.serialization.StringSerializer;

public class DataProducer {

    public static void main(String[] args) {
        Properties properties = new Properties();
        properties.put("bootstrap.servers", "192.168.8.172:9092");
        properties.put("key.serializer", StringSerializer.class.getName());
        properties.put("value.serializer", StringSerializer.class.getName());
        KafkaProducer producer = new KafkaProducer<String, String>(properties);

        produceEdges(producer);
    }

    private static void produceNodes(KafkaProducer<String, String> producer) {
        for (Person p : getPersons()) {
            ProducerRecord<String, String> record = new ProducerRecord<>(
                    "nebula", p.getId(), toJson(p));
            producer.send(record);
        }
        producer.flush();
    }

    private static void produceEdges(KafkaProducer<String, String> producer) {
        for (Follow f : getFollow()) {
            ProducerRecord<String, String> record = new ProducerRecord<>(
                    "nebula1", f.getId(), toJson(f));
            producer.send(record);
        }
        producer.flush();
    }


    private static String toJson(Person person) {
        JSONObject ob = (JSONObject) JSONObject.toJSON(person);
        return ob.toString();
    }

    private static String toJson(Follow follow) {
        JSONObject ob = (JSONObject) JSONObject.toJSON(follow);
        return ob.toString();
    }

    private static List<Person> getPersons() {
        List<Person> persons = new ArrayList<>();
        persons.add(new Person("100", "Tom", 12.0, true, 1.0));
        persons.add(new Person("101", "Tim", 13.0, false, 1.0));
        persons.add(new Person("102", "Nicole", 14.0, true, 1.0));
        persons.add(new Person("103", "Luna", 15.0, true, 1.0));
        persons.add(new Person("104", "Jena", 16.0, true, 1.0));
        persons.add(new Person("105", "Jena", 16.0, true, 1.0));
        persons.add(new Person("106", "Jena", 16.0, true, 1.0));
        persons.add(new Person("107", "Jena", 16.0, true, 1.0));
        persons.add(new Person("108", "Jena", 16.0, true, 1.0));
        return persons;
    }


    private static List<Follow> getFollow() {
        List<Follow> follows = new ArrayList<>();
        follows.add(new Follow("1", "101", "2", 100L, 1.0));
        follows.add(new Follow("2", "102", "19", 100L, 1.0));
        follows.add(new Follow("3", "103", "3", 100L, 1.0));
        follows.add(new Follow("4", "12", "1", 100L, 1.0));
        follows.add(new Follow("5", "13", "18", 100L, 1.0));
        follows.add(new Follow("6", "14", "17", 100L, 1.0));
        follows.add(new Follow("7", "15", "15", 100L, 1.0));

        return follows;
    }
}
