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

  const [recommendSongs, setRecommendSongs] = useState<Song[]>([])

  const fetchSongs = () => {
    fetch("http://localhost:8080/songs")
      .then(res => res.json())
      .then((data: Song[]) => setSongs(data));
  }

  useEffect(() => {
    fetchSongs()
  }, [])

  const handleSubmit = () => {
    fetch("http://localhost:8080/songs", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title, artist, tag }),
    })
    .then(res => {
      if (!res.ok) throw new Error(`API error: ${res.status}`);
      return res.json();
    })
    .then((newSong: Song) => {
      setSongs([...songs, newSong]);
      setTitle("");
      setArtist("");
      setTag("");
    })
    .catch(err => console.error("Failed to add song:", err));
  } 

  const fetchRecommend = () => {
    fetch("http://localhost:8080/recommend")
    .then(res => res.json())
    .then((data: Song[]) => setRecommendSongs(data));
  }

  return (
    <div>
      <h2>曲一覧</h2>
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
      <div>
        <h2>おすすめ</h2>
        <button onClick={fetchRecommend}>おすすめを見る</button>
          {recommendSongs.map(song => (
          <div key={song.id}>
            {song.title} - {song.artist} ({song.tag})
          </div>
        ))}
      </div>
    </div>
  )
}

export default App
