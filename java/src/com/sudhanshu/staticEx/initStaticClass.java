package src.com.sudhanshu.staticEx;

public class initStaticClass {
   static int a; 

   static {
        System.out.println("init");
        initStaticClass.a=4;
       
   }
   static void Print() {
        System.out.println("jsut chekc init");
       
   }
   public static void main(String[] args) {
        Print();


    
   }

}
