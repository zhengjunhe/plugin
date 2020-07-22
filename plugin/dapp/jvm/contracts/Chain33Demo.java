import com.sun.jna.*;
import java.util.*;
import java.lang.Long;

public class Chain33Demo {
   public interface Chain33Grpc extends Library {
        // GoSlice class maps to:
        // C type struct { void *data; GoInt len; GoInt cap; }
        public class GoSlice extends Structure {
            public static class ByValue extends GoSlice implements Structure.ByValue {}
            public Pointer data;
            public long len;
            public long cap;
            protected List getFieldOrder(){
                return Arrays.asList(new String[]{"data","len","cap"});
            }
        }

        // GoString class maps to:
        // C type struct { const char *p; GoInt n; }
        public class GoString extends Structure {
            public static class ByValue extends GoString implements Structure.ByValue {}
            public String p;
            public long n;
            protected List getFieldOrder(){
                return Arrays.asList(new String[]{"p","n"});
            }

        }

        // Foreign functions
//         public void StateDBRead(GoSlice.ByValue key, GoSlice.ByValue valueRead);
//         public void StateDBWrite(GoSlice.ByValue key, GoSlice.ByValue data);

        //ServerAddr, CaFile, ServerHostOverride string, tls bool
        public GoUint8 StartGrpcClient(GoString p0, GoString p1, GoString p2, GoUint8 p3);

        public void StopGrpcClient();

        //ExecFrozen(from string, amount int64)
        public GoUint8 ExecFrozen(GoString p0, GoInt64 p1);

        //ExecActive(from, execAddr string, amount int64) bool
        public GoUint8 ExecActive(GoString p0, GoString p1, GoInt64 p2);

        //ExecTransfer(from, to, execAddr string, amount int64)
        public GoUint8 ExecTransfer(GoString p0, GoString p1, GoString p2, GoInt64 p3);

        public GoSlice GetRandom();

        public GoString GetFrom();

        public GoUint8 SetState(GoSlice p0, GoSlice p1);

        public GoUint8 GetFromState(GoSlice p0, GoSlice p1);

        public GoInt32 GetValueSize(GoSlice p0);

        public GoUint8 SetLocalDB(GoSlice p0, GoSlice p1);

        public GoUint8 GetFromLocalDB(GoSlice p0, GoSlice p1);

        public GoInt32 GetLocalValueSize(GoSlice p0);
    }

   static public void main(String argv[]) {
        Chain33Grpc chain33Grpc = (Chain33Grpc) Native.load(
        "../contract/grpcClient.so", Chain33Grpc.class);

        //chain33 合约开发测试
        GoString from = chain33Grpc.GetFrom()
        System.out.printf("from is:%s\n" + from.p);

        chain33Grpc.ExecFrozen()

        // First, prepare data array
        byte[] key = new byte[]{104, 101, 108, 108, 111, 32, 99, 104, 97, 105, 110, 51, 51};
        Memory arr = new Memory(key.length);
         System.out.printf("key's length:%d\n", key.length);
        arr.write(0, key, 0, key.length);
        // fill in the GoSlice class for type mapping
        Blockchain.GoSlice.ByValue keySlice = new Blockchain.GoSlice.ByValue();
        keySlice.data = arr;
        keySlice.len = key.length;
        keySlice.cap = key.length;

        byte[] value = new byte[]{89, 101, 115, 44, 32, 71, 114, 101, 97, 116, 32, 67, 104, 105, 110, 97};
        System.out.printf("value's length:%d\n", value.length);
        Memory arr4Value = new Memory(value.length);
        arr4Value.write(0, value, 0, value.length);
        Blockchain.GoSlice.ByValue ValueSlice = new Blockchain.GoSlice.ByValue();
        ValueSlice.data = arr4Value;
        ValueSlice.len = value.length;
        ValueSlice.cap = value.length;
        dbOperation.StateDBWrite(keySlice, ValueSlice);

        //读取数据
        Blockchain.GoSlice.ByValue ValueSliceReadBack = new Blockchain.GoSlice.ByValue();
        Memory arr4ValueRead = new Memory(value.length);
        ValueSliceReadBack.data = arr4ValueRead;
        ValueSliceReadBack.len = value.length;
        ValueSliceReadBack.cap = value.length;
        dbOperation.StateDBRead(keySlice, ValueSliceReadBack);
        System.out.print("The data read back from chain33 is as below:\n");
        byte[] valueReadBack = ValueSliceReadBack.data.getByteArray(0, value.length);
        for(int i = 0; i < valueReadBack.length; i++){
            System.out.print(valueReadBack[i] + " ");
        }
        System.out.print("\n");




    }
}