package org.example;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.HashSet;
import java.util.Set;
import java.util.StringTokenizer;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapred.MapReduceBase;
import org.apache.hadoop.mapred.Mapper;
import org.apache.hadoop.mapred.OutputCollector;
import org.apache.hadoop.mapred.Reporter;

public class WordMapper extends MapReduceBase implements Mapper<Object, Text, CompositeKey, IntWritable>{
    private static final Set<String> stopWords = loadStopWords();

    private static Set<String> loadStopWords() {
        Configuration conf = new Configuration();

        try {
            FileSystem fs = FileSystem.get(conf);
            Set<String> readStopWords;
            try (BufferedReader br =
                         new BufferedReader(new InputStreamReader(fs.open(new Path("hdfs:/input/stopwords.txt"))))) {
                readStopWords = new HashSet<>();
                String line = br.readLine();
                while (line != null) {
                    readStopWords.add(line.toLowerCase());
                    line = br.readLine();
                }
            }
            return readStopWords;
        } catch (IOException e) {
            e.printStackTrace();
        }
        return new HashSet<>();
    }

    private final Text word = new Text();
    private final static IntWritable ONE = new IntWritable(1);

    @Override
    public void map(Object obj, Text value,
                    OutputCollector<CompositeKey, IntWritable> collector, Reporter report)
            throws IOException {
        StringTokenizer itr = new StringTokenizer(value.toString(), " ");

        while (itr.hasMoreTokens()) {
            word.set(itr.nextToken());
            if (!stopWords.contains(word.toString())) {
                collector.collect(new CompositeKey(word.toString(), 1), ONE);
            }
        }
    }

}
