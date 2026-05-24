package  src.com.sudhanshu.app;

import src.com.sudhanshu.staticEx.Human;

public class Main {
    public static void main(String[] args) {
        Human sid = new Human("aman", 22, false);
        Human harsh  =new  Human("Harsh", 22, false);
        System.out.println(sid);
        //cause it's link with the  object not the main blublprntwhic refer eeryon eth eocmomn place 
        System.out.println(sid.population);
        System.out.println(harsh);
        
    }    
}
