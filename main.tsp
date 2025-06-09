import "@typespec/http";
import "@typespec/protobuf";

using TypeSpec.Http;
using Protobuf;

@TypeSpec.service(#{ title: "Notes API" })
@server("http://localhost:3000", "dev")
@package({ name: "notesapi.v1",
        go_package:"github.com/the-gaurav-bhatt/notesapi" })
namespace NotesApi;

// --------------------------------------------------------------------------
// Messages
// --------------------------------------------------------------------------
@message model User        { @field(1) id: string;  @field(2) username: string; @field(3) password: string; }
@message model Credentials { @field(1) username: string; @field(2) password: string; }
@message model Note        { @field(1) id: string;  @field(2) title: string;    @field(3) content: string; @field(4) userId: string; }
@message model ListNotesResponse { @field(1) notes: Note[]; }
@message model Empty {}

// --------------------------------------------------------------------------
// Auth service – one parameter per RPC, so each gets @field(n)
// --------------------------------------------------------------------------
@route("/auth")
@Protobuf.service
interface Auth {
  @post @route("/register")
  register(@body  @field(1) user: User): User;

  @post @route("/login")
  login   (@body  @field(1) credentials: Credentials): User;
}

// --------------------------------------------------------------------------
// Notes service – path/body params all tagged
// --------------------------------------------------------------------------
@route("/notes")
@Protobuf.service
interface Notes {

  @get
  listNotes(): ListNotesResponse;

  @post
  createNote(@body @field(1) note: Note): Note;

  @get  @route("/{id}")
  getNote   (@path @field(1) id: string): Note;

  @put  @route("/{id}")
  updateNote(@path @field(1) id: string,
             @body @field(2) note: Note): Note;

  @delete @route("/{id}")
  deleteNote(@path @field(1) id: string): Empty;
}
