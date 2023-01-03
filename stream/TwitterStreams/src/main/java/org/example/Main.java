package org.example;

import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.*;
import org.apache.kafka.streams.kstream.*;

import java.time.Duration;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Properties;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

import static org.apache.kafka.streams.kstream.Suppressed.BufferConfig.unbounded;

public class Main {
    private static HashSet<String> BadWords = new HashSet<>() {
        {
            add("death");
            add("dead");
            add("crime");
            add("criminality");
            add("profane");
            add("profanity");
            add("violent");
            add("violence");
            add("drugs");
            add("abuse");
            add("nudity");
            add("sex");
            add("prison");
            add("rape");
            add("fraud");
            add("destroy");
            add("kill");
            add("poison");
            add("virus");
            add("disease");
            add("danger");
            add("dangerous");
            add("corruption");
            add("corrupt");
        }};
    public static void main(String[] args) throws Exception {
        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, "streams-pipe");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "broker:9092");
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.String().getClass());
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass());

        Duration windowSize = Duration.ofSeconds(5);
        Duration advanceSize = Duration.ofSeconds(5);

        TimeWindows tumblingWindow = TimeWindows.of(windowSize).advanceBy(advanceSize).grace(Duration.ZERO);

        final StreamsBuilder builder = new StreamsBuilder();
        builder.stream("twitter-streams-in", Consumed.with(Serdes.String(), Serdes.String()))
                .groupByKey()
                .windowedBy(tumblingWindow)
                .aggregate(
                        () -> "0",
                        (key, value, total) -> {
                            String[] lines = value.split("\n");
                            ArrayList<String> categories = new ArrayList<>();
                            for (String line: lines) {
                                categories.addAll(Arrays.stream(line.split("~\\|~")).toList());
                            }
                            ArrayList<String> words = new ArrayList<>();
                            for (String category: categories) {
                                words.addAll(Arrays.stream(category.split(",")).toList());
                            }
                            AtomicInteger subtotal = new AtomicInteger();
                            words.forEach(x -> {
                                if (BadWords.contains(x)) {
                                    subtotal.getAndIncrement();
                                }
                            });
                            var intTotal = Integer.parseInt(total);
                            return String.valueOf(intTotal + subtotal.get());
                        },
                        Materialized.with(Serdes.String(), Serdes.String()))
                .suppress(Suppressed.untilWindowCloses(unbounded()))
                .toStream()
                .map((wk, value) -> KeyValue.pair(wk.key(), value + "===" + Instant.now().getEpochSecond()))
                .peek((key, value) -> System.out.println("key " + key + " value " + value))
                .to("twitter-streams-out", Produced.with(Serdes.String(), Serdes.String()));

        final Topology topology = builder.build();

        final KafkaStreams streams = new KafkaStreams(topology, props);
        final CountDownLatch latch = new CountDownLatch(1);

        // attach shutdown handler to catch control-c
        Runtime.getRuntime().addShutdownHook(new Thread("streams-shutdown-hook") {
            @Override
            public void run() {
                streams.close();
                latch.countDown();
            }
        });

        try {
            streams.start();
            latch.await();
        } catch (Throwable e) {
            System.exit(1);
        }
        System.exit(0);
    }
}