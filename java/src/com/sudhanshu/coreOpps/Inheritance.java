package src.com.sudhanshu.coreOpps;

public class Inheritance {
    static class A{
        String Name;
        public A(String name){
            this.Name=name;
        }



    }
    static class B extends A{
        public B(String Name){
            super(Name);
        }
        @Override
        public String toString(){
            return  this.Name;
        }

    }
    public static void main(String[] args) {
        System.out.println(new B("sid"));
        
    }

}
