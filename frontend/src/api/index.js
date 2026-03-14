import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  withCredentials: true,
})

// Posts
export const getPosts = (params) => api.get('/posts', { params })
export const getPost = (id) => api.get(`/posts/${id}`)
export const createPost = (data) => api.post('/posts', data)
export const updatePost = (id, data) => api.put(`/posts/${id}`, data)
export const deletePost = (id) => api.delete(`/posts/${id}`)

// Categories & Tags
export const getCategories = () => api.get('/categories')
export const getCategory = (name) => api.get(`/categories/${name}`)
export const getTags = () => api.get('/tags')
export const getTag = (name) => api.get(`/tags/${name}`)
export const getArchives = () => api.get('/archives')

// Moments
export const getMoments = (params) => api.get('/moments', { params })
export const createMoment = (data) => api.post('/moments', data)
export const deleteMoment = (id) => api.delete(`/moments/${id}`)

// Comments
export const getComments = (params) => api.get('/comments', { params })
export const createComment = (data) => api.post('/comments', data)
export const deleteComment = (id) => api.delete(`/comments/${id}`)

// Likes
export const addLike = (data) => api.post('/likes', data)
export const removeLike = (data) => api.delete('/likes', { data })

// Auth
export const getMe = () => api.get('/auth/me')
export const logout = () => api.post('/auth/logout')

// Playground
export const runCode = (code) => api.post('/playground/run', { code })

// Upload
export const uploadFile = (formData) => api.post('/upload', formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
})
