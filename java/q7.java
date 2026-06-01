public class q7 {
    static class PaymentService {
        private class CreditCard implements PaymentMethod {
            public void process(double amount) {
                System.out.println("Processing payment of " + amount);
            }
        }
        public  PaymentMethod getPaymentMethod(){
            return  new CreditCard();
            
        }
    }
    interface PaymentMethod {
        void process(double amount);
    }

    public static void main(String[] args) {
        
        PaymentService ps = new PaymentService();
        ps.getPaymentMethod().process(100.0);
    }
}
