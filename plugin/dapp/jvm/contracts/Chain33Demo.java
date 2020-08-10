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

        public Boolean StartDefaultGrpcClient(GoString.ByValue p0);
        //ServerAddr, CaFile, ServerHostOverride string, tls bool
        public Boolean StartGrpcClient(GoString.ByValue p0, GoString.ByValue p1, GoString.ByValue p2, Boolean p3);

        public void StopGrpcClient();

        //ExecFrozen(from string, amount int64)
        public Boolean ExecFrozen(GoString p0, long p1);

        //ExecActive(from, execAddr string, amount int64) bool
        public Boolean ExecActive(GoString.ByValue p0, GoString.ByValue p1, Long p2);

        //ExecTransfer(from, to, execAddr string, amount int64)
        public Boolean ExecTransfer(GoString.ByValue p0, GoString.ByValue p1, GoString.ByValue p2, Long p3);

        public GoSlice GetRandom();

        public GoString GetFrom();

        public Boolean SetState(GoSlice.ByValue p0, GoSlice.ByValue p1);

        public Boolean GetFromState(GoSlice.ByValue p0, GoSlice.ByValue p1);

        public Integer GetValueSize(GoSlice.ByValue p0);

        public Boolean SetLocalDB(GoSlice.ByValue p0, GoSlice.ByValue p1);

        public Boolean GetFromLocalDB(GoSlice.ByValue p0, GoSlice.ByValue p1);

        public Integer GetLocalValueSize(GoSlice.ByValue p0);
    }

   static public void main(String[] args) {
        System.out.println("The input parameter is as below:");
        for (String arg : args) {
            System.out.println(arg);
        }
        System.out.println("-----------------------------");

        Chain33Grpc chain33Grpc = (Chain33Grpc) Native.load(
        "/root/contract/grpcClient.so", Chain33Grpc.class);

       Chain33Grpc.GoString.ByValue serverAddr = new Chain33Grpc.GoString.ByValue();
       serverAddr.p = "127.0.0.1:8802";
       serverAddr.n = serverAddr.p.length();
       if (!chain33Grpc.StartDefaultGrpcClient(serverAddr)) {
           System.out.println("Failed to StartDefaultGrpcClient");
           return;
       }

        //chain33 合约开发测试
        Chain33Grpc.GoString from = chain33Grpc.GetFrom();
        System.out.printf("from is:%s\n" + from.p);

        long amount = 100;
        Boolean execResult;
        execResult = chain33Grpc.ExecFrozen(from, amount);
        if (!execResult) {
            System.out.println("Failed to ExecFrozen");
        }
        System.out.println("Succeed to ExecFrozen");


        // First, prepare data array
//         byte[] key = new byte[]{104, 101, 108, 108, 111, 32, 99, 104, 97, 105, 110, 51, 51};
//         Memory arr = new Memory(key.length);
//          System.out.printf("key's length:%d\n", key.length);
//         arr.write(0, key, 0, key.length);
//         // fill in the GoSlice class for type mapping
//         Chain33Grpc.GoSlice.ByValue keySlice = new Chain33Grpc.GoSlice.ByValue();
//         keySlice.data = arr;
//         keySlice.len = key.length;
//         keySlice.cap = key.length;
//
//         byte[] value = new byte[]{89, 101, 115, 44, 32, 71, 114, 101, 97, 116, 32, 67, 104, 105, 110, 97};
//         System.out.printf("value's length:%d\n", value.length);
//         Memory arr4Value = new Memory(value.length);
//         arr4Value.write(0, value, 0, value.length);
//         Chain33Grpc.GoSlice.ByValue ValueSlice = new Chain33Grpc.GoSlice.ByValue();
//         ValueSlice.data = arr4Value;
//         ValueSlice.len = value.length;
//         ValueSlice.cap = value.length;
//         dbOperation.StateDBWrite(keySlice, ValueSlice);
//
//         //读取数据
//         Chain33Grpc.GoSlice.ByValue ValueSliceReadBack = new Chain33Grpc.GoSlice.ByValue();
//         Memory arr4ValueRead = new Memory(value.length);
//         ValueSliceReadBack.data = arr4ValueRead;
//         ValueSliceReadBack.len = value.length;
//         ValueSliceReadBack.cap = value.length;
//         dbOperation.StateDBRead(keySlice, ValueSliceReadBack);
//         System.out.print("The data read back from chain33 is as below:\n");
//         byte[] valueReadBack = ValueSliceReadBack.data.getByteArray(0, value.length);
//         for(int i = 0; i < valueReadBack.length; i++){
//             System.out.print(valueReadBack[i] + " ");
//         }
//         System.out.print("\n");
    }
}