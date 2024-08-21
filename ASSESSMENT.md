Homework Assignment - GoLang

Overview :

You are part of an organisation that owns a Bike Sharing Platform that allows its customers to rent and return bikes to bike docking stations. Customers tap their access card given by the platform which contains a user UUID on the bike docking station to access available bikes. A backend assignment service determines which bike to assign to the customer and unlocks the bike.
You own the backend that allows users to assign and un-assign the bikes from bike docking station. For simplicity, customers already registered on the platform rent bikes from the docking station on their way to work, and return the bike at the end of the day to the same docking station.

Key Considerations

    Assignment operation should choose the least used bike;
    The same bike cannot be assigned to several other users;
    When bike is returned (unassign operation) it is unavailable to be assigned for the next 5 minutes;
    A bike should be auto unassigned from customer after being in use by a single user for more than 24 hours in a row.

Technologies to use

    GoLang 1.17
    Framework : Any, standard libraries are highly. For web frameworks : chi and mux packages
    Container : Docker
    Relational DB (MS SQL, MySQL, PostgreSQL)
    use a migration library to handle database migrations

Task: Create a RESTful service providing API’s which will be used by the docking station to assign and un-assign bike to users.

    Create a Docker Compose Configuration
    Store DB Credentials in a Config File
    Create an HTTP Server
    Create Tables to Store Users, Bikes, and Their Mapping
    come up with fields that are needed and could be needed for current and future usages
    Create Initial Migrations with Users and Bikes
    (3 people), bikes table (2 pieces)
    Create REST Endpoints for the following
        assigns a bike to a user
        unassigns a bike
        returns bikes available for assignment

Additional Considerations:

    Think about possible errors and handle them accordingly;
    Use logging tools to log important blocks of the workflows;
    Write unit tests to check the logic.

Additional questions (Interview only)

    Platforms Customer Service Team user could have roles (Admin, Supervisor), default is Customer
    Add an Admin user using a migration or update an existing user to be an Admin
    Add a check that a bike cannot be assigned to an Admin. Implement an error for this case (Bad Request)
    Add a test or modify an existing one to check the logic
