# Song Library API

- **List song data with filtering and pagination:**
    ```http
    GET /songs
    ```
     ```http
    GET /songs?group=Muse&page=2
    ```
    - queries for filtering:
        - group
        - song
        - releaseDate
            - Release date field can only be either in format "DD.MM.YYYY", "MM.YYYY" or "YYYY" 
        - text
        - link
    - queries for pagination:
        - page
        - pageSize
    - sample output:
    ```json
    [{
        "id": 11,
        "group": "Muse",
        "song": "Supermassive Black Hole",
        "releaseDate": "16.07.2006",
        "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
        "totalVerses": 6
    }]
    ```
- **Get song's text data**
    - required parameter: `id`
    ```http
    GET /songs/:id/text
    ```
    - queries:
        - verse
            - If verse is set to 0 (default), then return the whole text, otherwise return a paginated text (a specified verse)
    - sample output:
      - `?verse=0`
        
        ```json
        {
            "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
        }
        ```
      - `?verse=1`
  
        ```json
        {
            "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?"
        }
        ```
- **Adding new song data**
    ```http
    POST /songs
    ```
    - input body:
    ```json
    {
    "group": "Muse",
    "song": "Supermassive Black Hole"
    }
    ```
- **Update song info:**
    - required parameter: `id`
     ```http
    PATCH /song/:id
    ```
    - input body:
    ```json
    {
        "group": "Muse",
        "song": "Uprising",
        "releaseDate": "04.08.2009",
        "link": "https://www.youtube.com/watch?v=w8KQmps-Sog",
        "text": [
            "(Come on)",
            "Paranoia is in bloom\nThe PR transimissions will resume\nThey'll try to push drugs that keep us all dumbed down"
        ]
    }
    ```
- **Delete song info:**
    - required parameter: `id`
     ```http
    DELETE /song/:id
    ```
---
### Start
**Make sure there is an .env file. Create it from the example** `.env.example` **file**

- install dependencies:
```
go get
```
- run project:
```
go run main.go 
```