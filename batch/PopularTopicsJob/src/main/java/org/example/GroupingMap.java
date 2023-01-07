package org.example;

import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapred.MapReduceBase;
import org.apache.hadoop.mapred.Mapper;
import org.apache.hadoop.mapred.OutputCollector;
import org.apache.hadoop.mapred.Reporter;

import java.io.IOException;

public class GroupingMap extends MapReduceBase implements Mapper<Object, Text, CompositeKey, IntWritable>{

    @Override
    public void map(Object arg0, Text item,
                    OutputCollector<CompositeKey, IntWritable> collector, Reporter r)
            throws IOException {
        String str= item.toString();
        String[] items = str.split("\t");
        CompositeKey ck= new CompositeKey(items[0], Integer.parseInt(items[items.length -1]));
        collector.collect(ck, new IntWritable(ck.getCount()));
    }

}
