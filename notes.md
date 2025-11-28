 Input Types (Request body for mutations):what data the user must send during mutations.
input PaginationInput {
  skip: Int
  take: Int
}
Used for pagination:
skip → how many items to skip
take → how many items to return(Like LIMIT/OFFSET)
