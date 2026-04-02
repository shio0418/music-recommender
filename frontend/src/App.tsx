import { useEffect, useState } from 'react'

type Song = {
  id: number
  title: string
  artist: string
  tag: string
  youtubeUrl?: string
  videoId?: string
  thumbnailUrl?: string
  embedUrl?: string
  viewCount?: number
  likeCount?: number
  commentCount?: number
  score?: number
}

function App() {
  const [songs, setSongs] = useState<Song[]>([])

  const [title, setTitle] = useState('')
  const [artist, setArtist] = useState('')
  const [tag, setTag] = useState('')
  const [youtubeUrl, setYoutubeUrl] = useState('')

  const [recommendSongs, setRecommendSongs] = useState<Song[]>([])

  const fetchSongs = () => {
    fetch('http://localhost:8080/songs')
      .then(res => res.json())
      .then((data: Song[]) => setSongs(data))
      .catch(err => console.error('Failed to fetch songs:', err))
  }

  useEffect(() => {
    fetchSongs()
  }, [])

  const handleSubmit = () => {
    const payload = youtubeUrl.trim()
      ? { youtubeUrl: youtubeUrl.trim(), tag: tag.trim() }
      : { title: title.trim(), artist: artist.trim(), tag: tag.trim() }

    fetch('http://localhost:8080/songs', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
    .then(res => {
      if (!res.ok) throw new Error(`API error: ${res.status}`)
      return res.json()
    })
    .then((newSong: Song) => {
      setSongs(prev => [...prev, newSong])
      setTitle('')
      setArtist('')
      setTag('')
      setYoutubeUrl('')
    })
    .catch(err => console.error('Failed to add song:', err))
  }

  const fetchRecommend = () => {
    fetch('http://localhost:8080/recommend')
      .then(res => res.json())
      .then((data: Song[]) => setRecommendSongs(data))
      .catch(err => console.error('Failed to fetch recommendations:', err))
  }

  return (
    <div>
      <h2>曲一覧</h2>
      {songs.map(song => (
        <div key={song.id}>
          <div>{song.title} - {song.artist} ({song.tag})</div>
          {typeof song.score === 'number' && (
            <div>score: {song.score.toFixed(2)}</div>
          )}
          {song.embedUrl && (
            <iframe
              title={`player-${song.id}`}
              width="300"
              height="169"
              src={song.embedUrl}
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
              referrerPolicy="strict-origin-when-cross-origin"
              allowFullScreen
            />
          )}
        </div>
      ))}
      <div>
        <h2>曲を追加</h2>

        <input
          placeholder="YouTube URL（入れると自動取得）"
          value={youtubeUrl}
          onChange={(e) => setYoutubeUrl(e.target.value)}
        />
        
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
