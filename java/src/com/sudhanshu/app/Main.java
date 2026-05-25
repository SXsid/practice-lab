package  src.com.sudhanshu.app;

import src.com.sudhanshu.staticEx.Human;

public class Main {
    // main is the first this to run when creating a class but 
    //it it's not static we have to create a boject first which is not possible
    //hece frist / main is awlsy static it's not depne on object but on clas it slef
    public static void main(String[] args) {

        Human sid = new Human("aman", 22, false);
        Human harsh  =new  Human("Harsh", 22, false);
        System.out.println(sid);
        //cause it's link with the  object not the main blublprntwhic refer eeryon eth eocmomn place 
        System.out.println(sid.population);
        System.out.println(harsh);
        
    }    
}
