package src.com.sudhanshu.coreOpps;

public class Inheritance {
    static class A{
        int Age=22;
        String Name;
        public A(String name){
            this.Name=name;
        }



    }
    static class B extends A{
        String city;
        public B(String Name, String city){
            super(Name);
            this.city=city;
        }
        @Override
        public String toString(){
            return  this.Name + this.Age;
        }

    }
    public static void main(String[] args) {
        System.out.println(new B("sid","pali"));
        
    }

}
