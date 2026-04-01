import { useEffect, useState } from 'react'

type Song = {
    id: number;
    title: string;
    artist: string;
    tag: string;
}

function App() {
  const [songs, setSongs] = useState<Song[]>([])

  const [title, setTitle] = useState("");
  const [artist, setArtist] = useState("");
  const [tag, setTag] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/songs")
      .then(res => res.json())
      .then((data: Song[]) => setSongs(data));
  }, [])

  const handleSubmit = () => {
    fetch("http://localhost:8080/songs", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title, artist, tag }),
    })
  } 

  return (
    <div>
      <h1>曲一覧</h1>
      {songs.map(song => (
        <div key={song.id}>
          {song.title} - {song.artist} ({song.tag})
        </div>
      ))}
      <div>
        <h2>曲を追加</h2>
        
        <input
          placeholder="タイトル"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />

        <input
          placeholder="アーティスト"
          value={artist}
          onChange={(e) => setArtist(e.target.value)}
        />

        <input
          placeholder="タグ"
          value={tag}
          onChange={(e) => setTag(e.target.value)}
        />

        <button onClick={handleSubmit}>追加</button>

      </div>
    </div>
  )
}

export default App
