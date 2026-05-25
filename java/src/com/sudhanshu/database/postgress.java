package src.com.sudhanshu.database;

//singeton  should be applied only one object
public class postgress {
    String DSN;
    //no one can create a object out of it ousteide this file mean only one instace of it will created;
    private  postgress(String DSN){
        this.DSN=DSN;

    }
    private static postgress instance;
    //expose the publi funtion which sotre the single insae of the ojbect
    public static postgress getInstace(String dsn){
        if (postgress.instance==null){
            postgress.instance= new postgress(dsn);
        }

        return  postgress.instance;
    }
    
}

class Main{
    public static void main(String[] args) {
        postgress ints1 = postgress.getInstace("dsn1");
        postgress ints2 = postgress.getInstace("dsn2");
        System.out.println(ints1.DSN);
        System.out.println(ints2.DSN);
        
    }

}
