import axios from 'axios'

const baseURL = 'http://localhost:5500/api/v1'

export const useApi = () => {
  const api = axios.create({
    baseURL: baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  return {
    async getAllRouters() {
      try {
        const { data } = await api.get(`${baseURL}/routers`)
        console.log(data);
        
        return data
      } catch (error) {
        console.error('API Error:', error)
        throw new Error('Failed to fetch routers')
      }
    },

    async createRouter(router: { name: string; ip_address: string }) {
      try {
        const { data } = await api.post(`${baseURL}/routers`, router)
        return data
      } catch (error) {
        console.error('API Error:', error)
        throw new Error('Failed to create router')
      }
    },

    async connectRouter(data: { router_from_ip: string; router_to_ip: string }) {
      try {
        const { data: responseData } = await api.post(`${baseURL}/routers/connection`, {
          router_from_ip: data.router_from_ip,
          router_to_ip: data.router_to_ip
        })
        return responseData
      } catch (error) {
        console.error('API Error:', error)
        throw new Error('Failed to connect router')
      }
    },

    async configureRouter(data: {
      routerId: number;
      name: string;
      status: string;
    }) {
      try {
        const { data: responseData } = await api.patch(`${baseURL}/routers/configure`, {
          routerId: data.routerId,
          name: data.name,
          status: data.status,
        })
        return responseData
      } catch (error) {
        console.error('API Error:', error)
        throw new Error('Failed to configure router')
      }
    },

    async configureRouterPort(data: {
      routerId: number;
      portNumber: number;
      protocol: string;
      status: string;
      speed: string;
      duplexMode: string;
      description?: string;
    }) {
      try {
        const { data: responseData } = await api.patch(`${baseURL}/ports/configure`, {
          routerId: data.routerId.toString(),
          portNumber: data.portNumber,
          protocol: data.protocol,
          status: data.status,
          speed: data.speed,
          duplexMode: data.duplexMode,
          description: data.description,
        })
        return responseData
      } catch (error) {
        console.error('API Error:', error)
        throw new Error('Failed to configure router')
      }
    },

    async getRouterConnections() {
      try {
        const { data } = await api.get(`${baseURL}/routers/connections`)
        return data
      } catch (error) {
        console.error('API Error:', error)
        throw new Error('Failed to fetch router connections')
      }
    }
  }
} 