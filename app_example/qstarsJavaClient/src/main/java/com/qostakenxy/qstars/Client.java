package com.qostakenxy.qstars;

import com.sun.jna.*;
import java.util.*;
import java.lang.Long;

public class Client {
	public interface Awesome extends Library {

		// GoSlice class maps to:
		// C type struct { void *data; GoInt len; GoInt cap; }
		public class GoSlice extends Structure {
			public static class ByValue extends GoSlice implements Structure.ByValue {
			}

			public Pointer data;
			public long len;
			public long cap;

			protected List getFieldOrder() {
				return Arrays.asList(new String[] { "data", "len", "cap" });
			}
		}

		// GoString class maps to:
		// C type struct { const char *p; GoInt n; }
		public class GoString extends Structure {
			public static class ByValue extends GoString implements Structure.ByValue {
			}

			public String p;
			public long n;

			protected List getFieldOrder() {
				return Arrays.asList(new String[] { "p", "n" });
			}

		}

		// Foreign functions
		public long Add(long a, long b);

		public double Cosine(double val);

		public void Sort(GoSlice.ByValue vals);

		public long Log(GoString.ByValue str);
	}

	public interface Aninterface extends Library {

		// GoSlice class maps to:
		// C type struct { void *data; GoInt len; GoInt cap; }
		public class GoString extends Structure {
			public static class ByValue extends GoString implements Structure.ByValue {
			}

			public String p;
			public long n;

			protected List getFieldOrder() {
				return Arrays.asList(new String[] { "p", "n" });
			}

		}

		// Foreign functions
		public String AccountCreate();
		public String QSCQueryAccount(String url);
		public Integer QSCKVStoreSet(String k, String v, String privkey, String chain);
		public String QSCKVStoreGet(String url);
		public String QSCtransfer(String ul,String a,String privkey,String chain,String ac,String seq,String g);

	}

	static public void main(String argv[]) {
		Aninterface aninterface = (Aninterface) Native.loadLibrary("./Aninterface.so", Aninterface.class);
		aninterface.AccountCreate();
		String out = aninterface.QSCQueryAccount("cosmosaccaddr1nskqcg35k8du3ydhntkcqjxtk254qv8me943mv");
		Integer outkvresult = aninterface.QSCKVStoreSet("8", "Hyert", "9Rg9mNEXVh9aUsxJ74Ogqe8O6wrBw8EeMhyK/GgHcfUsGprPgC7YXH6YEwGM+eXmc7oV1ci7ivlxo7k6amd3Lg", "test-chain-kkbpwi");
		String kvout = aninterface.QSCKVStoreGet("7");
		String transferout = aninterface.QSCtransfer("cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9", "3QSC1", "8O2FbbnWpIff/cs5anJK13+RwNdh0GO7PmDhDIXEHKkbVsM3LPqu319fbBGg3j1ocY3xqhxra8oPEHjMbddVeA", "test-chain-kkbpwi", "2", "1", "1");
		System.out.println(out);
		System.out.println(outkvresult);
		System.out.println(kvout);
		System.out.println(transferout);
		
	}

	static public void mainExample(String argv[]) {

		Awesome awesome = (Awesome) Native.loadLibrary("./Aninterface.so", Awesome.class);

		System.out.printf("awesome.Add(12, 99) = %s\n", awesome.Add(12, 99));
		System.out.printf("awesome.Cosine(1.0) = %s\n", awesome.Cosine(1.0));

		// Call Sort
		// First, prepare data array
		long[] nums = new long[] { 53, 11, 5, 2, 88 };
		Memory arr = new Memory(nums.length * Native.getNativeSize(Long.TYPE));
		arr.write(0, nums, 0, nums.length);
		// fill in the GoSlice class for type mapping
		Awesome.GoSlice.ByValue slice = new Awesome.GoSlice.ByValue();
		slice.data = arr;
		slice.len = nums.length;
		slice.cap = nums.length;
		awesome.Sort(slice);
		System.out.print("awesome.Sort(53,11,5,2,88) = [");
		long[] sorted = slice.data.getLongArray(0, nums.length);
		for (int i = 0; i < sorted.length; i++) {
			System.out.print(sorted[i] + " ");
		}
		System.out.println("]");

		// Call Log
		Awesome.GoString.ByValue str = new Awesome.GoString.ByValue();
		str.p = "Hello Java!";
		str.n = str.p.length();
		System.out.printf("msgid %d\n", awesome.Log(str));

	}
}