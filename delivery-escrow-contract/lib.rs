#![cfg_attr(not(feature = "std"), no_std, no_main)]

#[ink::contract]
mod delivery_escrow {
    
    #[ink(storage)]
    pub struct DeliveryEscrow {
        is_delivered: bool,
        balance: Balance,
        buyer: AccountId,
        transporter: AccountId,
        proof_of_delivery: Option<Hash>,
        delivery_timestamp: Option<Timestamp>,
        timeout_duration: Timestamp,
        dispute: bool,
    }

    impl DeliveryEscrow {
        // This constructor intializes the delivery status to false
        #[ink(constructor)]
        pub fn new(buyer: AccountId, transporter: AccountId) -> Self {
            Self {
                is_delivered: false,
                balance: 0,
                buyer,
                transporter,
                delivery_timestamp: None,
                dispute: false,
                proof_of_delivery: None,
                timeout_duration: 0,
            }
        }

        // Mark when the product is delivered
        #[ink(message)]
        pub fn mark_delivered(&mut self) -> bool {
            self.is_delivered = true;
            self.is_delivered
        }

        // Return whether the product has been delivered
        #[ink(message)]
        pub fn is_delivered(&self) -> bool {
            self.is_delivered
        }

        // Accept money from the caller and store amount in balance field
        #[ink(message, payable)]
        pub fn deposit(&mut self) {
            let amount = self.env().transferred_value();
            assert!(amount > 0, "You must deposit more than 0");

            self.balance = self.balance.checked_add(amount).expect("Overflow occurred");
        }

        // This function  checks if the product was delivered
        // Transfers the escrowed funds to the seller account
        // Resets the balance to zero
        #[ink(message)]
        pub fn release(&mut self, seller: AccountId) {
            assert!(self.is_delivered, "Delivery not confirmed");

            let amount = self.balance;
            assert!(amount > 0, "No funds to release");

            // Calculate the split amounts
            let transporter_share = amount.checked_mul(10).expect("Overflow in calculation") / 100;
            let seller_share = amount.checked_sub(transporter_share).expect("Overflow in calculation");

            self.balance = 0;

            // Transfer 10% to the transporter
            let result_transporter = self.env().transfer(self.transporter, transporter_share);
            assert!(result_transporter.is_ok(), "Transporter transfer failed");

            // Transfer 90% to the seller
            let result_seller = self.env().transfer(seller, seller_share);
            assert!(result_seller.is_ok(), "Seller transfer failed");
        }


        
        // This is incase the seller fails to confirm delivery
        #[ink(message)]
        pub fn refund(&mut self) {
            // Only the buyer can initiate a refund
            let caller = self.env().caller();
            assert_eq!(caller, self.buyer, "Only buyer can request refund");

            assert!(!self.is_delivered, "Cannot refund after delivery confirmed");

            let amount = self.balance;
            assert!(amount > 0, "No funds to refund");

            self.balance = 0;

            let result = self.env().transfer(self.buyer, amount);
            assert!(result.is_ok(), "Refund transfer failed");
        }

        // When the seller delivers the goods you need to show proof
        #[ink(message)]
        pub fn submit_proof(&mut self, proof: Hash) {
            let caller = self.env().caller();

            // Prevent buyer from submitting proof
            assert_ne!(caller, self.buyer, "Buyer cannot submit proof");

            // Only allow if not already delivered
            assert!(!self.is_delivered, "Already marked as delivered");

            self.proof_of_delivery = Some(proof);
            self.delivery_timestamp = Some(self.env().block_timestamp());
            self.is_delivered = true;
        }

        // The function for the buyer to confirm delivery.
        #[ink(message)]
        pub fn confirm_delivery(&mut self) {
            let caller = self.env().caller();
            assert_eq!(caller, self.buyer, "Only buyer can confirm delivery");

            // Must be a proof from the seller first
            assert!(
                self.proof_of_delivery.is_some(),
                "No proof submitted by seller"
            );

            self.is_delivered = true;
            self.dispute = false; // Resolved
        }

        //Raise a dispute if the buyer doesnt agree
        #[ink(message)]
        pub fn raise_dispute(&mut self) {
            let caller = self.env().caller();
            assert_eq!(caller, self.buyer, "Only buyer can raise dispute");

            assert!(
                self.is_delivered,
                "Delivery has not been marked by seller"
            );

            self.dispute = true;
        }

    }

    /// Unit tests in Rust are normally defined within such a `#[cfg(test)]`
    /// module and test functions are marked with a `#[test]` attribute.
    /// The below code is technically just normal Rust code.
    #[cfg(test)]
    mod tests {
        use crate::delivery_escrow::DeliveryEscrow;
        // use ink_env::{self, DefaultEnvironment};
        use ink_env::hash::{Blake2x256, HashOutput};
        type Hash = <DefaultEnvironment as ink_env::Environment>::Hash;


        #[ink::test]
        fn confirm_delivery_works() {
            let accounts = ink::env::test::default_accounts::<ink::env::DefaultEnvironment>();

            // Instantiate the contract
            let mut contract = DeliveryEscrow::new(accounts.alice, accounts.charlie);

            // Set caller as buyer
            ink_env::test::set_caller::<ink_env::DefaultEnvironment>(accounts.alice);

            // Call deposit with 1000 balance and seller as Bob
            ink_env::test::set_value_transferred::<ink_env::DefaultEnvironment>(1000);
            contract.deposit();

            // Confirm delivery (caller is still Alice)
            contract.confirm_delivery();

            // Assert that delivery is marked true
            assert_eq!(contract.is_delivered, true);
        }

        #[ink::test]
        fn test_submit_proof_works() {
            let accounts = ink_env::test::default_accounts::<ink_env::DefaultEnvironment>();
            let mut contract = DeliveryEscrow {
                is_delivered: false,
                balance: 0,
                buyer: accounts.bob,
                transporter: accounts.charlie, // <-- add this line
                proof_of_delivery: None,
                delivery_timestamp: None,
                timeout_duration: 10000,
                dispute: false,
            };

            // Set the caller to the seller (not the buyer)
            ink_env::test::set_caller::<ink_env::DefaultEnvironment>(accounts.alice);

            // Example of data to hash
            let delivery_data = b"Order #123 delivered to Bob";

            // a buffer to hold the hash
            let mut output = <Blake2x256 as HashOutput>::Type::default();

            // Compute the hash
            ink::env::hash_bytes::<Blake2x256>(delivery_data, &mut output);
            contract.submit_proof(output);

            assert_eq!(contract.proof_of_delivery, Some(output));
            assert_eq!(contract.is_delivered, true);
            assert!(contract.delivery_timestamp.is_some());
        }

    }
}
