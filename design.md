# Design

## Determine Target Audience

MKB/Small business and individuals.
Full internet bureau's with a little more technological knowledge. (They will want to scale their front-end and have LB failover)

## Use cases

MoSCoW: **M**ust, **S**hould, **C**ould, **W**ould
Where **M**ust is used, you can consider it to be in the MVP.

| MoSCoW | User Story |
|--------|------------|
| M | As any user I must be able to login securely |
| M | As a client I must be able to browse through your hosting products |
| M | As a client I must be able to configure hosting products |
| M | As a client I must be able to add hosting configurations to cart |
| M | As a client I must be able to also buy a domain name |
| M | As a client I must be able to see statistics of my hosting products |
| M | As a client I must be able to upgrade my hosting products |
| M | As a client I must be able to see all invoices |
| M | As a client I must be able to download an invoices |
| M | As a client I want to get a hosting configuration recommendation before buying |
| M | As a client I want to get a new password when I forget it |
| S | As a client I want have an option to pay each quarter |
| S | As a client I want have an option to pay each year |
| S | As a client I want to have scheduled backups of my files each X period |
| S | As a client I want to have an on-demand backup of my files |
| S | As a client I want to have scheduled backups of my database each X period |
| M | As a client I want to have an on-demand backup of my database |
| M | As a client I want be able to access my files remotely |
| M | As a client I want be able to access my database remotely |
| ? | As a client I want to be able to SSH into my (data)container |
| ? | As a client I want to have a custom IP per website/company |
|  | As an account manager I must be able to add a new client |
|  | As an account manager I must be able to create an invoice |
|  | As an account manager I must be able to buy a domain for a client |
|  | As an account manager I must be able to configure a hosting product for a client |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |
|  | As an account manager I  |

## Minimum Viable Product

There is a marketing site where clients can browse products and read about FAQ and other marketing related topics. Terms of Service etc.
An ordering system is in place to pick and configure the hosting product.
And a checkout can be done and payments are possible via through some payment provider. client can reach clear documentation on how to manage her (Wordpress) website.

### Skateboard version
client can buy a hosted wordpress service. One containerized wordpress with period backups of both files and database with some limit yet to be defined.
Container cluster needs to scale horizontally automatically.

Should I use k8s for this, or use plain VPS'es and scale them up myself? Scaling down means you are losing clients.

For the domain address we can use transip domain only.
I'm not planning on being an email provider, maybe it can be outsourced to *transip/gmail/another hosting provider* that provides email as a service.
Can clients plugin their own email server/service? Will they need DNS access? Or do they own DNS access?

### Technical details
For an MVP the following moving parts are required:
- IP address
- Load Balancer
- 1 Controller/Backend Node (hosting the "api" and "control-panel")
- 1 Front-end Node (hosting the "portal")
- N Hosting Nodes

#### API
- [ ] Golang/NodeJS/Python/PHP?
- [ ] Authentication
- [ ] Authorization
- [ ] Backend API
- [ ] Shopping cart API
- [ ] Hosting API

#### Control Panel
Front-end for end authenticated end users that bought a hosting product. They can manage everything from this control panel.
- [ ] Manage paid services
- [ ] Buy more paid services

#### Portal
Front-end for potential clients. They can get information and sign up and buy products.
- [ ] Services overview
- [ ] Shopping cart for hosting products
- [ ] Marketing stuff that makes potential clients wanna BUY IT NOW
- [ ] FAQ
- [ ] Terms of Service
- [ ] Contact page
- [ ] Support page


