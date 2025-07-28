import express from 'express'
import cors from 'cors'
import helmet from 'helmet'
import morgan from 'morgan'
import dotenv from 'dotenv'

// Load environment variables
dotenv.config()

const app = express()
const port = process.env.PORT || 3000

// Middleware
app.use(helmet())
app.use(cors())
app.use(morgan('combined'))
app.use(express.json())
app.use(express.urlencoded({ extended: true }))

// Routes
app.get('/', (req, res) => {
  res.json({
    message: 'Welcome to {{.ProjectName}} API',
    owner: '{{.Owner}}',
    timestamp: new Date().toISOString()
  })
})

app.get('/api/status', (req, res) => {
  res.json({
    status: 'OK',
    environment: process.env.NODE_ENV || 'development',
    timestamp: new Date().toISOString()
  })
})

// Error handling middleware
app.use((err: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
  console.error(err.stack)
  res.status(500).json({
    error: 'Something went wrong!',
    message: process.env.NODE_ENV === 'development' ? err.message : 'Internal server error'
  })
})

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({
    error: 'Route not found',
    path: req.originalUrl
  })
})

app.listen(port, () => {
  console.log(`🚀 {{.ProjectName}} server running on port ${port}`)
  console.log(`📝 Owner: {{.Owner}}`)
  console.log(`🌍 Environment: ${process.env.NODE_ENV || 'development'}`)
}) 