package src.com.sudhanshu.staticEx;

public class innerClass {

    class Inner{
        String name;
        Inner(String name){
            this.name=name;
        }


    }
    static class InnerStaticcanbeusedinpsvm{
        String name;
        public InnerStaticcanbeusedinpsvm(String name){
            this.name=name;
        }
        
    }
    public static void main(String[] args) {
        InnerStaticcanbeusedinpsvm a = new InnerStaticcanbeusedinpsvm("amna");
        System.out.println(a.name);

    }
    public void run(){
        Inner a = new Inner("aman");
        Inner b = new Inner("harsh");
        System.out.println(a.name);
        System.out.println(b.name);
    }
    
}
