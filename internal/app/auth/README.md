# auth module

This is a separate auth module in modular monolith fashion.

This is loosely coupled in ride-sharing-golang-api project and won't be used tightly coupled to ride-sharing business.

Instead, a creation of customer or driver would create a customer or driver in trip module with a reference of that user from auth module.
