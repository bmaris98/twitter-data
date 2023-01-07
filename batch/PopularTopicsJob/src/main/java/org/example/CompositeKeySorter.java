package org.example;

import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableComparator;

public class CompositeKeySorter extends WritableComparator {
    public CompositeKeySorter(){
        super(CompositeKey.class, true);
    }

    @Override
    public int compare(WritableComparable a, WritableComparable b) {
        CompositeKey ck1 =(CompositeKey)a;
        CompositeKey ck2 =(CompositeKey)b;

        return ck1.getWord().compareTo(ck2.getWord());
    }

}
