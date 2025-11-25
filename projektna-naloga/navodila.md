# Projektna naloga: Razpravljalnica
<p align="center">
  <img src="logo.png">
</p>

Ustvarite distribuirano spletno storitev **Razpravljalnica**, ki bo namenjena izmenjavi mnenj med uporabniki o različnih temah. 

Razpravljalnici se lahko pridružijo novi **uporabniki**, ki nato vanjo dodajajo nove **teme** in znotraj posameznih tem objavljajo **sporočila**. Razpravljalnica omogoča, da se lahko uporabnik na eno ali več tem naroči in sproti prejema iz strežnika sporočila, ki se objavljajo znotraj naročenih tem. Podprto je tudi **všečkanje** sporočil. Za vsako objavljeno sporočilo se beleži število prejetih všečkov. Uporabniki lahko lastna sporočila tudi **urejajo/brišejo**. 

Storitev (strežnik) napišite v programskem jeziku Go. Uporablja naj ogrodje **gRPC** za komunikacijo z odjemalci (uporabniki). Prav tako napišite odjemalca, ki bo znal komunicirati s strežnikom in bo podpiral vse operacije, ki jih ponuja Razpravljalnica. Za komunikacijo znotraj storitve (med strežniki) lahko uporabite poljubno rešitev (rpc). 

## Podrobnosti naloge

Projekto nalogo rešujete v parih. Na spletni učilnici najdete [anketo](https://ucilnica.fri.uni-lj.si/mod/choicegroup/view.php?id=60665) preko katere se razdelite v pare.
Vašo rešitev v obliki modula objavite v javnem repozitoriju GitHub/Gitlab in preko spletne učilnice [oddate](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=16663) povezavo do njega. Vaš repozitorij naj bo do izteka roka za oddajo zaseben, potem pa ga naredite javnega. **Rok za oddajo projektne naloge je 12. 1. 2026**. Po izteku roka bomo organizirali zagovore na katerih boste predstavili vaše rešitve.

Podano imate spodnjo specifikacijo v jeziku Protocol Buffer, ki definira programski vmesnik Razpravljalnice, kot je viden iz stališča odjemalca (uporabnika):

```protobuf
syntax = "proto3";

package razpravljalnica;

option go_package = "/razpravljalnica";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";
////////////////////////////////////////////////////////////////////////////////
// Basic data types
////////////////////////////////////////////////////////////////////////////////

message User {
  int64 id = 1;
  string name = 2;
}

message Topic {
  int64 id = 1;
  string name = 2;
}

message Message {
  int64 id = 1;
  int64 topic_id = 2;
  int64 user_id = 3;
  string text = 4;
  google.protobuf.Timestamp created_at = 5;
  int32 likes = 6;
}

message Like {
  int64 topic_id = 1;
  int64 message_id = 2;
  int64 user_id = 3; // user who liked the message
}

enum OpType {
  OP_POST   = 0; // add a message to a topic
  OP_LIKE   = 1; // like a message
  OP_DELETE = 2; // delete a message
  OP_UPDATE = 3; // update a message
}

message NodeInfo {
  string node_id = 1;
  string address = 2;
}

////////////////////////////////////////////////////////////////////////////////
// Data plane
////////////////////////////////////////////////////////////////////////////////

service  MessageBoard{
  // Writes (only to head)

  // Creates a new user and assigns it an id
  rpc CreateUser(CreateUserRequest) returns (User);

  // Creates a new topic to which users can post messages
  rpc CreateTopic(CreateTopicRequest) returns (Topic);

  // Post a message to a topic; Succed only if the User and the Topic exist in the data base.
  rpc PostMessage(PostMessageRequest) returns (Message);

  // Update an existing message. Allowed only for the user who posted the message.
  rpc UpdateMessage(UpdateMessageRequest) returns (Message);

  // Delete an existing message. Allowed only for the user who posted the message.
  rpc DeleteMessage(DeleteMessageRequest) returns (google.protobuf.Empty);

  // Like an existing message. Return the message with the new number of likes.
  rpc LikeMessage(LikeMessageRequest) returns (Message);

  // Request a node to which a subscription can be opened.
  rpc GetSubcscriptionNode(SubscriptionNodeRequest) returns (SubscriptionNodeResponse);

  // Reads go to the tail node

  // Returns all the topics
  rpc ListTopics(google.protobuf.Empty) returns (ListTopicsResponse);

  // Returns messages in a topic
  rpc GetMessages(GetMessagesRequest) returns (GetMessagesResponse);

  // Subscribe to topics; goes to the node returned by head
  rpc SubscribeTopic(SubscribeTopicRequest) returns (stream MessageEvent);

}

message CreateUserRequest {
  string name = 1;
}

message CreateTopicRequest {
  string name = 1;
}

message PostMessageRequest {
  int64 topic_id = 1;
  int64 user_id = 2;
  string text = 3;
}

message LikeMessageRequest {
  int64 message_id = 1;
  int64 user_id = 2; // user who posted the like
}

message ListTopicsResponse {
  repeated Topic topics = 1;
}

message GetMessagesRequest {
  int64 topic_id = 1;
  int64 from_message_id = 2; // starting id of the message (0 from beggining)
  int32 limit = 3; // max number of messages
}
message GetMessagesResponse {
  repeated Message messages = 1;
}

message SubscribeTopicRequest {
  repeated int64 topic_id = 1;
  int64 user_id = 2;
  int64 from_message_id = 3; // starting id of the message
  string subscribe_token = 4; // token generated by the head used to authorize the subscription
}

message SubscriptionNodeRequest {
  int64 user_id = 1;
  repeated int64 topic_id = 2;
}

message SubscriptionNodeResponse {
  string subscribe_token = 1; // subscription token to be presented to the returned node
  NodeInfo node = 2;
}

message MessageEvent {
  int64 sequence_number = 1; // monotonically increasing event number
  OpType op = 2; // type of event
  Message message = 3;
  google.protobuf.Timestamp event_at = 4; // timestamp of the event
}

////////////////////////////////////////////////////////////////////////////////
// Control plane
////////////////////////////////////////////////////////////////////////////////

// Return the the head and the tail node address
service ControlPlane {
  rpc GetClusterState(google.protobuf.Empty) returns (GetClusterStateResponse);
}

message GetClusterStateResponse {
  NodeInfo head = 1;
  NodeInfo tail = 2;
}
```

## Kriterij ocenjevanja

Primarna naloga je, da implementirate zgornji programski vmesnik strežnika in pripravite odjemalca, ki bo s strežnikom znal komunicirati in uporabljati vse oddaljene operacije, ki jih strežnik ponuja. Strežnik naj vse podatke hrani v delovnem pomnilniku (*ang. in-memory storage*) Pripravite tudi demonstracijo delovanja strežnika in odjemalca.

 - **Ocena 6-7:** Strežnik je samo en in omogoča hkratno delo z večimi odjemalci. Implementiran mora biti celotni programski vmesnik po zgornji specifikaciji.

 - **Ocena 7-8:** Uporabite [verižno replikacijo](https://www.cs.cornell.edu/home/rvr/papers/OSDI04.pdf), zato da porazdelite bralne dostope med več vozlišč. Predpostavite, da so vozlišča popolnoma zanesljiva. Vsi pisalni dostopi se izvajajo nad glavo. Naročnine na teme pa glava po nekem ključu (uravnoteženje obremenitve) porazdeli med vozlišča. Enkratni bralni dostopi se izvedejo nad repom verige. Programski vmesnik ustrezno razširite tako, da bo omogočal delovanje verižne replikacije. Uporabite lahko poljubno tehnologijo za izvedbo oddaljenih klicev.

 - **Ocena 9-10:** Uporabite verižno replikacijo, zato da porazdelite bralne dostope med več vozlišč. Predpostavite, da vozlišča niso zanesljiva in lahko pride do odpovedi vozlišča v verigi. V takem primeru se mora veriga ponovno vzpostaviti in zapisi uskladiti. Vsi pisalni dostopi se izvajajo nad glavo. Naročnine na teme pa glava po nekem ključu (uravnoteženje obremenitve) porazdeli med vozlišča. Enkratni bralni dostopi se izvedejo nad repom verige. V verigo naj bo možno dodati nova vozlišča in jih odvzemati. Za preverjanje dostopnosti vozlišč v verigi, dodajanje novih in prevezovanje v primeru odpovedi implementirajte tudi nadzorno ravnino (dodatni strežnik). Ustrezno razširite programski vmesnik. 

 - **Bonus:**
    -  Uporabite pakete za pisanje ukaznih vmesnikov ukazne lupine, kot so [cli](https://cli.urfave.org/), [cobra](https://pkg.go.dev/github.com/spf13/cobra), [kong](https://github.com/alecthomas/kong).
    - Uporabite pakete za pisanje grafičnih vmesnikov za strežnik in odjemalec npr. [tview](https://github.com/rivo/tview).
    - Go ima vgrajeno [podporo](https://go.dev/doc/tutorial/fuzz) za pisanje testov. Pripravite nabor testov, ki validirajo delovanje vaše Razpravljalnice.
    - Predvidite, da lahko pride do odpovedi nadzorne ravnine. Uvedite več vozlišč v nadzorno ravnino in uporabite protokol [raft](https://repozitorij.uni-lj.si/Dokument.php?id=215152&lang=slv), da poskrbite za primere odpovedi. Poslužite se lahko katere izmed obstoječih, na primer [HashiCorp](https://github.com/hashicorp/raft) ali [diploma FRI](https://github.com/Timcek/raft) protokola raft v Go.