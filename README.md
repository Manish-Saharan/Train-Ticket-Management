# Train-Ticket-Management

Train Ticket Management is a simple command-line application that allows users to purchase tickets for a train journey. The application is built using Golang and gRPC for communication between the client and server.

## Features

1. **Purchase Ticket:** 
Users can submit a purchase request for a train ticket by providing details such as origin, destination, user information (first name, last name, email), and the price paid. The user is allocated a seat in either section A or section B.

2. **View Receipt:** 
Retrieve the details of a ticket purchase receipt by providing the user's email address.

3. **View Users By Section:** 
View the list of users in a specific train section (A or B).

4. **Remove User:** 
Remove a user from the train by providing their email address.

5. **Modify Seat:** 
Modify a user's seat by providing their details and the desired train section.

