<template>
  <div id="app" class="min-h-screen bg-gray-100">
    <!-- Header -->
    <header class="bg-blue-600 text-white shadow-lg">
      <div class="container mx-auto px-4 py-6">
        <h1 class="text-3xl font-bold text-center">🌐 Domain Check</h1>
        <p class="text-center mt-2 text-blue-100">Domain availability and information checker</p>
      </div>
    </header>

    <!-- Main Content -->
    <main class="container mx-auto px-4 py-8">
      <!-- Domain Check Form -->
      <div class="max-w-2xl mx-auto mb-8">
        <div class="bg-white rounded-lg shadow-md p-6">
          <h2 class="text-xl font-semibold mb-4">Domain Availability Check</h2>

          <!-- Tab Selection -->
          <div class="flex mb-4 border-b">
            <button
              @click="activeTab = 'single'"
              :class="['px-4 py-2 font-medium', activeTab === 'single' ? 'border-b-2 border-blue-500 text-blue-600' : 'text-gray-600']"
            >
              Single Domain
            </button>
            <button
              @click="activeTab = 'all-extensions'"
              :class="['px-4 py-2 font-medium', activeTab === 'all-extensions' ? 'border-b-2 border-blue-500 text-blue-600' : 'text-gray-600']"
            >
              Check All Extensions
            </button>
          </div>

          <!-- Single Domain Check -->
          <div v-if="activeTab === 'single'" class="space-y-4">
            <div class="flex gap-4">
              <input
                v-model="domainInput"
                type="text"
                placeholder="Enter full domain (e.g., google.com)"
                class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                @keyup.enter="checkDomain"
              />
              <button
                @click="checkDomain"
                :disabled="loading || !domainInput.trim()"
                class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ loading ? 'Checking...' : 'Check' }}
              </button>
            </div>
          </div>

          <!-- All Extensions Check -->
          <div v-if="activeTab === 'all-extensions'" class="space-y-4">
            <div class="flex gap-4">
              <input
                v-model="domainNameInput"
                type="text"
                placeholder="Enter domain name only (e.g., metehansaral)"
                class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                @keyup.enter="checkAllExtensions"
              />
              <button
                @click="checkAllExtensions"
                :disabled="bulkLoading || !domainNameInput.trim()"
                class="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ bulkLoading ? 'Checking All...' : 'Check All Extensions' }}
              </button>
            </div>
            <p class="text-sm text-gray-600">
              🚀 This will check your domain name with all {{ totalExtensions }}+ available extensions (.com, .net, .org, etc.)
            </p>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="mt-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {{ error }}
          </div>
        </div>
      </div>

      <!-- Current Check Result -->
      <div v-if="currentResult" class="max-w-2xl mx-auto mb-8">
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold mb-4">Check Result</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-gray-600">Domain</label>
              <p class="text-lg font-semibold">{{ currentResult.name }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-600">Status</label>
              <p class="text-lg font-semibold" :class="getStatusClass(currentResult.status)">
                {{ getStatusText(currentResult.status) }}
              </p>
            </div>
            <div v-if="currentResult.ip">
              <label class="text-sm font-medium text-gray-600">IP Address</label>
              <p class="text-lg font-mono">{{ currentResult.ip }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-600">Checked At</label>
              <p class="text-sm text-gray-500">{{ formatDate(currentResult.checked_at) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Bulk Check Results -->
      <div v-if="bulkResult" class="max-w-6xl mx-auto mb-8">
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex justify-between items-center mb-6">
            <h3 class="text-xl font-semibold">🚀 Domain Availability Results for "{{ bulkResult.domain_name }}"</h3>
            <span class="text-sm text-gray-500">{{ bulkResult.total_time_ms }}ms</span>
          </div>

          <!-- Summary Cards -->
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <div class="bg-green-50 border border-green-200 rounded-lg p-4 text-center">
              <div class="text-2xl font-bold text-green-600">{{ bulkResult.available_count }}</div>
              <div class="text-sm text-green-700">Available</div>
            </div>
            <div class="bg-red-50 border border-red-200 rounded-lg p-4 text-center">
              <div class="text-2xl font-bold text-red-600">{{ bulkResult.unavailable_count }}</div>
              <div class="text-sm text-red-700">Taken</div>
            </div>
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 text-center">
              <div class="text-2xl font-bold text-yellow-600">{{ bulkResult.error_count }}</div>
              <div class="text-sm text-yellow-700">Errors</div>
            </div>
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 text-center">
              <div class="text-2xl font-bold text-blue-600">{{ bulkResult.total_extensions }}</div>
              <div class="text-sm text-blue-700">Total Checked</div>
            </div>
          </div>

          <!-- Recommended Domains -->
          <div v-if="bulkResult.summary.recommended_domains.length > 0" class="mb-6">
            <h4 class="text-lg font-semibold mb-3 text-green-600">🎯 Recommended Available Domains</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              <div v-for="domain in bulkResult.summary.recommended_domains.slice(0, 6)" :key="domain"
                   class="bg-green-50 border border-green-200 rounded-lg p-3 text-center">
                <span class="font-semibold text-green-800">{{ domain }}</span>
                <div class="text-xs text-green-600 mt-1">Available!</div>
              </div>
            </div>
          </div>

          <!-- Tabs for detailed results -->
          <div class="border-b mb-4">
            <nav class="flex space-x-4">
              <button @click="resultTab = 'available'"
                      :class="['py-2 px-4 font-medium', resultTab === 'available' ? 'border-b-2 border-green-500 text-green-600' : 'text-gray-600']">
                Available ({{ bulkResult.available_count }})
              </button>
              <button @click="resultTab = 'taken'"
                      :class="['py-2 px-4 font-medium', resultTab === 'taken' ? 'border-b-2 border-red-500 text-red-600' : 'text-gray-600']">
                Taken ({{ bulkResult.unavailable_count }})
              </button>
              <button @click="resultTab = 'alternatives'"
                      :class="['py-2 px-4 font-medium', resultTab === 'alternatives' ? 'border-b-2 border-blue-500 text-blue-600' : 'text-gray-600']">
                Alternatives
              </button>
            </nav>
          </div>

          <!-- Available Domains -->
          <div v-if="resultTab === 'available'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            <div v-for="result in bulkResult.available_domains" :key="result.domain.name"
                 class="bg-green-50 border border-green-200 rounded-lg p-3">
              <div class="font-semibold text-green-800">{{ result.domain.name }}</div>
              <div class="text-xs text-green-600">{{ result.domain.response_time_ms }}ms</div>
            </div>
          </div>

          <!-- Taken Domains -->
          <div v-if="resultTab === 'taken'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            <div v-for="result in bulkResult.unavailable_domains" :key="result.domain.name"
                 class="bg-red-50 border border-red-200 rounded-lg p-3">
              <div class="font-semibold text-red-800">{{ result.domain.name }}</div>
              <div class="text-xs text-red-600">{{ result.domain.ip || 'Registered' }}</div>
            </div>
          </div>

          <!-- Alternative Suggestions -->
          <div v-if="resultTab === 'alternatives'" class="space-y-4">
            <div v-if="bulkResult.summary.alternative_suggestions.length > 0">
              <h5 class="font-semibold text-blue-600 mb-2">💡 Alternative Suggestions</h5>
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                <div v-for="suggestion in bulkResult.summary.alternative_suggestions" :key="suggestion"
                     class="bg-blue-50 border border-blue-200 rounded-lg p-3 text-center">
                  <span class="font-semibold text-blue-800">{{ suggestion }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Domain History -->
      <div class="max-w-4xl mx-auto">
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold">Check History</h3>
            <button
              @click="loadDomains"
              class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
            >
              Refresh
            </button>
          </div>

          <div v-if="domains.length === 0" class="text-center py-8 text-gray-500">
            No domains checked yet. Try checking a domain above!
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full table-auto">
              <thead>
                <tr class="bg-gray-50">
                  <th class="px-4 py-2 text-left">Domain</th>
                  <th class="px-4 py-2 text-left">Status</th>
                  <th class="px-4 py-2 text-left">IP Address</th>
                  <th class="px-4 py-2 text-left">Checked At</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="domain in domains" :key="domain.name + domain.checked_at" class="border-t">
                  <td class="px-4 py-2 font-medium">{{ domain.name }}</td>
                  <td class="px-4 py-2">
                    <span :class="getStatusClass(domain.status)">
                      {{ getStatusText(domain.status) }}
                    </span>
                  </td>
                  <td class="px-4 py-2 font-mono text-sm">{{ domain.ip || '-' }}</td>
                  <td class="px-4 py-2 text-sm text-gray-500">{{ formatDate(domain.checked_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'App',
  data() {
    return {
      activeTab: 'single',
      resultTab: 'available',
      domainInput: '',
      domainNameInput: '',
      loading: false,
      bulkLoading: false,
      error: null,
      currentResult: null,
      bulkResult: null,
      domains: [],
      totalExtensions: 228
    }
  },
  mounted() {
    this.loadDomains()
  },
  methods: {
    async checkDomain() {
      if (!this.domainInput.trim()) return

      this.loading = true
      this.error = null
      this.currentResult = null
      this.bulkResult = null

      try {
        const response = await axios.post('/api/check-domain', {
          domain: this.domainInput.trim()
        })

        if (response.data.success) {
          this.currentResult = response.data.data.domain
          this.domainInput = ''
          // Listeyi yenile
          this.loadDomains()
        } else {
          this.error = response.data.message || 'Check failed'
        }
      } catch (error) {
        this.error = error.response?.data?.message || 'Network error occurred'
      } finally {
        this.loading = false
      }
    },

    async checkAllExtensions() {
      if (!this.domainNameInput.trim()) return

      this.bulkLoading = true
      this.error = null
      this.currentResult = null
      this.bulkResult = null

      try {
        const response = await axios.post('/api/check-all-extensions', {
          domain_name: this.domainNameInput.trim()
        })

        if (response.data.success) {
          this.bulkResult = response.data.data
          this.resultTab = 'available'
          this.domainNameInput = ''
          // Listeyi yenile
          this.loadDomains()
        } else {
          this.error = response.data.message || 'Bulk check failed'
        }
      } catch (error) {
        this.error = error.response?.data?.message || 'Network error occurred'
      } finally {
        this.bulkLoading = false
      }
    },

    async loadDomains() {
      try {
        const response = await axios.get('/api/domains')
        if (response.data.success) {
          this.domains = response.data.data || []
          // En son kontrol edilenleri başta göster
          this.domains.reverse()
        }
      } catch (error) {
        console.error('Failed to load domains:', error)
      }
    },

    formatDate(dateString) {
      const date = new Date(dateString)
      return date.toLocaleString()
    },

    getStatusClass(status) {
      switch(status) {
        case 'Available':
          return 'text-green-600'
        case 'Registered':
          return 'text-red-600'
        case 'Error':
          return 'text-yellow-600'
        default:
          return 'text-gray-600'
      }
    },

    getStatusText(status) {
      switch(status) {
        case 'Available':
          return '✅ Available'
        case 'Registered':
          return '❌ Registered'
        case 'Error':
          return '⚠️ Error'
        default:
          return status || 'Unknown'
      }
    }
  }
}
</script>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}
</style>
