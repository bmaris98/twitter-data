package org.example;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableUtils;

public class CompositeKey implements WritableComparable {
    private String word;
    private int count;

    public String getWord() {
        return word;
    }



    public int getCount() {
        return count;
    }

    public void setCount(int count) {
        this.count = count;
    }

    public CompositeKey() {

    }

    public CompositeKey(String word, int count) {

        this.word = word;
        this.count = count;
    }

    @Override
    public void readFields(DataInput in) throws IOException {
        word = WritableUtils.readString(in);
        count = WritableUtils.readVInt(in);
    }

    @Override
    public void write(DataOutput out) throws IOException {
        WritableUtils.writeString(out, word);
        WritableUtils.writeVInt(out,count);
    }

    @Override
    public int compareTo(Object o) {
        CompositeKey ck = (CompositeKey)o;
        return ck.word.compareTo(word);
    }

    @Override
    public int hashCode() {
        return 23*word.hashCode();
    }

    @Override
    public boolean equals(Object obj) {
        CompositeKey ck = (CompositeKey)obj;
        return ck.word.equals(word);
    }

    @Override
    public String toString() {
        return word;
    }

}
