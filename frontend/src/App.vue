<template>
  <div id="app" :class="darkMode ? 'dark' : ''" class="min-h-screen transition-colors duration-300">
    <!-- Dark mode toggle -->
    <div class="fixed top-4 right-4 z-50">
      <button
        @click="toggleDarkMode"
        class="p-3 rounded-full bg-white dark:bg-gray-800 shadow-lg hover:shadow-xl transition-all duration-300"
        :title="darkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
      >
        <svg v-if="darkMode" class="w-6 h-6 text-yellow-500" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z" clip-rule="evenodd"></path>
        </svg>
        <svg v-else class="w-6 h-6 text-gray-700" fill="currentColor" viewBox="0 0 20 20">
          <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z"></path>
        </svg>
      </button>
    </div>

    <!-- Header -->
    <header class="bg-gradient-to-r from-blue-600 to-purple-600 dark:from-gray-800 dark:to-gray-900 text-white shadow-lg">
      <div class="container mx-auto px-4 py-8">
        <h1 class="text-4xl font-bold text-center mb-2">🌐 Domain Check</h1>
        <p class="text-center text-blue-100 dark:text-gray-300">Modern domain availability checker with real-time updates</p>
      </div>
    </header>

    <!-- Main Content -->
    <main class="container mx-auto px-4 py-8 bg-gray-50 dark:bg-gray-900 min-h-screen">
      <!-- Domain Check Form -->
      <div class="max-w-3xl mx-auto mb-8">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 border border-gray-200 dark:border-gray-700">
          <h2 class="text-2xl font-bold mb-6 text-gray-800 dark:text-white">Domain Availability Check</h2>

          <!-- Tab Selection -->
          <div class="flex mb-6 border-b border-gray-200 dark:border-gray-700">
            <button
              @click="activeTab = 'single'"
              :class="[
                'px-6 py-3 font-medium transition-all duration-200 rounded-t-lg',
                activeTab === 'single'
                  ? 'border-b-2 border-blue-500 text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200'
              ]"
            >
              Single Domain
            </button>
            <button
              @click="activeTab = 'all-extensions'"
              :class="[
                'px-6 py-3 font-medium transition-all duration-200 rounded-t-lg',
                activeTab === 'all-extensions'
                  ? 'border-b-2 border-blue-500 text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200'
              ]"
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
                class="flex-1 px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
                @keyup.enter="checkDomain"
              />
              <button
                @click="checkDomain"
                :disabled="loading || !domainInput.trim()"
                class="px-8 py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white rounded-lg font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg hover:shadow-xl"
              >
                <span v-if="loading" class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Checking...
                </span>
                <span v-else>Check</span>
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
                class="flex-1 px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
                @keyup.enter="checkAllExtensionsWebSocket"
              />
              <button
                @click="checkAllExtensionsWebSocket"
                :disabled="wsLoading || !domainNameInput.trim()"
                class="px-8 py-3 bg-gradient-to-r from-green-600 to-green-700 hover:from-green-700 hover:to-green-800 text-white rounded-lg font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg hover:shadow-xl"
              >
                <span v-if="wsLoading" class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Checking...
                </span>
                <span v-else>Check All Extensions</span>
              </button>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              🚀 Real-time checking with WebSocket for instant updates
            </p>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="mt-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 rounded-lg">
            {{ error }}
          </div>
        </div>
      </div>

      <!-- WebSocket Progress -->
      <div v-if="wsProgress && !wsProgress.isComplete" class="max-w-3xl mx-auto mb-8">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 border border-gray-200 dark:border-gray-700">
          <h3 class="text-xl font-bold mb-6 text-gray-800 dark:text-white">Real-time Progress</h3>

          <!-- Progress Bar -->
          <div class="mb-6">
            <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mb-3">
              <span>Progress: {{ wsProgress.checkedCount }} / {{ wsProgress.totalExtensions }}</span>
              <span>{{ Math.round((wsProgress.checkedCount / wsProgress.totalExtensions) * 100) }}%</span>
            </div>
            <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3">
              <div
                class="bg-gradient-to-r from-blue-500 to-purple-500 h-3 rounded-full transition-all duration-300"
                :style="{ width: (wsProgress.checkedCount / wsProgress.totalExtensions) * 100 + '%' }"
              ></div>
            </div>
          </div>

          <!-- Current Domain -->
          <div v-if="wsProgress.currentDomain" class="mb-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div class="flex justify-between items-center">
              <span class="font-medium text-gray-700 dark:text-gray-300">Currently checking:</span>
              <span class="font-mono text-blue-800 dark:text-blue-300">{{ wsProgress.currentDomain.domain }}</span>
            </div>
            <div class="flex justify-between items-center mt-2">
              <span class="text-sm text-gray-600 dark:text-gray-400">Status:</span>
              <span :class="getStatusClass(wsProgress.currentDomain.status)">
                {{ getStatusText(wsProgress.currentDomain.status) }}
              </span>
            </div>
          </div>

          <!-- Summary -->
          <div class="grid grid-cols-3 gap-4 text-center">
            <div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
              <div class="text-2xl font-bold text-green-600 dark:text-green-400">{{ wsProgress.availableCount }}</div>
              <div class="text-sm text-green-700 dark:text-green-300">Available</div>
            </div>
            <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
              <div class="text-2xl font-bold text-red-600 dark:text-red-400">{{ wsProgress.unavailableCount }}</div>
              <div class="text-sm text-red-700 dark:text-red-300">Taken</div>
            </div>
            <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
              <div class="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{{ wsProgress.errorCount }}</div>
              <div class="text-sm text-yellow-700 dark:text-yellow-300">Errors</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Current Check Result -->
      <div v-if="currentResult" class="max-w-3xl mx-auto mb-8">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 border border-gray-200 dark:border-gray-700">
          <h3 class="text-xl font-bold mb-6 text-gray-800 dark:text-white">Check Result</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="text-sm font-medium text-gray-600 dark:text-gray-400">Domain</label>
              <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ currentResult.name }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-600 dark:text-gray-400">Status</label>
              <p class="text-lg font-semibold" :class="getStatusClass(currentResult.status)">
                {{ getStatusText(currentResult.status) }}
              </p>
            </div>
            <div v-if="currentResult.ip">
              <label class="text-sm font-medium text-gray-600 dark:text-gray-400">IP Address</label>
              <p class="text-lg font-mono text-gray-900 dark:text-white">{{ currentResult.ip }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-600 dark:text-gray-400">Checked At</label>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ formatDate(currentResult.checked_at) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- WebSocket Results -->
      <div v-if="wsProgress && wsProgress.isComplete" class="max-w-6xl mx-auto mb-8">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 border border-gray-200 dark:border-gray-700">
          <div class="flex justify-between items-center mb-8">
            <h3 class="text-2xl font-bold text-gray-800 dark:text-white">🚀 Results for "{{ wsProgress.domainName }}"</h3>
            <span class="text-sm text-gray-500 dark:text-gray-400">{{ wsProgress.totalTime }}ms</span>
          </div>

          <!-- Summary Cards -->
          <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-6 text-center">
              <div class="text-3xl font-bold text-green-600 dark:text-green-400">{{ wsProgress.availableCount }}</div>
              <div class="text-sm text-green-700 dark:text-green-300">Available</div>
            </div>
            <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6 text-center">
              <div class="text-3xl font-bold text-red-600 dark:text-red-400">{{ wsProgress.unavailableCount }}</div>
              <div class="text-sm text-red-700 dark:text-red-300">Taken</div>
            </div>
            <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-6 text-center">
              <div class="text-3xl font-bold text-yellow-600 dark:text-yellow-400">{{ wsProgress.errorCount }}</div>
              <div class="text-sm text-yellow-700 dark:text-yellow-300">Errors</div>
            </div>
            <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6 text-center">
              <div class="text-3xl font-bold text-blue-600 dark:text-blue-400">{{ wsProgress.totalExtensions }}</div>
              <div class="text-sm text-blue-700 dark:text-blue-300">Total Checked</div>
            </div>
          </div>

          <!-- Available Domains -->
          <div v-if="wsProgress.availableDomains.length > 0" class="mb-8">
            <h4 class="text-xl font-bold mb-4 text-green-600 dark:text-green-400">🎯 Available Domains</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div v-for="domain in wsProgress.availableDomains.slice(0, 12)" :key="domain.domain"
                   class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4 hover:shadow-md transition-shadow">
                <div class="font-semibold text-green-800 dark:text-green-300">{{ domain.domain }}</div>
                <div class="text-xs text-green-600 dark:text-green-400 mt-1">{{ domain.responseTime }}ms</div>
              </div>
            </div>
          </div>

          <!-- Taken Domains (limited display) -->
          <div v-if="wsProgress.unavailableDomains.length > 0" class="mb-8">
            <h4 class="text-xl font-bold mb-4 text-red-600 dark:text-red-400">❌ Taken Domains (showing first 10)</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div v-for="domain in wsProgress.unavailableDomains.slice(0, 10)" :key="domain.domain"
                   class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 hover:shadow-md transition-shadow">
                <div class="font-semibold text-red-800 dark:text-red-300">{{ domain.domain }}</div>
                <div class="text-xs text-red-600 dark:text-red-400 mt-1">{{ domain.ip || 'Registered' }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Domain History -->
      <div class="max-w-6xl mx-auto">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 border border-gray-200 dark:border-gray-700">
          <div class="flex justify-between items-center mb-6">
            <h3 class="text-xl font-bold text-gray-800 dark:text-white">Check History</h3>
            <button
              @click="loadDomains"
              class="px-4 py-2 bg-gray-600 dark:bg-gray-700 text-white rounded-lg hover:bg-gray-700 dark:hover:bg-gray-600 transition-colors"
            >
              Refresh
            </button>
          </div>

          <div v-if="domains.length === 0" class="text-center py-12 text-gray-500 dark:text-gray-400">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <p class="mt-2">No domains checked yet. Try checking a domain above!</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-gray-200 dark:border-gray-700">
                  <th class="px-4 py-3 text-left text-sm font-medium text-gray-600 dark:text-gray-400">Domain</th>
                  <th class="px-4 py-3 text-left text-sm font-medium text-gray-600 dark:text-gray-400">Status</th>
                  <th class="px-4 py-3 text-left text-sm font-medium text-gray-600 dark:text-gray-400">IP Address</th>
                  <th class="px-4 py-3 text-left text-sm font-medium text-gray-600 dark:text-gray-400">Checked At</th>
                  <th class="px-4 py-3 text-left text-sm font-medium text-gray-600 dark:text-gray-400">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="domain in domains" :key="domain.name + domain.checked_at" class="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                  <td class="px-4 py-3 font-medium text-gray-900 dark:text-white">{{ domain.name }}</td>
                  <td class="px-4 py-3">
                    <span :class="getStatusClass(domain.status)">
                      {{ getStatusText(domain.status) }}
                    </span>
                  </td>
                  <td class="px-4 py-3 font-mono text-sm text-gray-900 dark:text-white">{{ domain.ip || '-' }}</td>
                  <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ formatDate(domain.checked_at) }}</td>
                  <td class="px-4 py-3">
                    <button
                      @click="showWhoisInfo(domain.name)"
                      class="px-3 py-1 bg-blue-500 text-white text-sm rounded-lg hover:bg-blue-600 transition-colors"
                    >
                      WHOIS
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Whois Modal -->
      <div v-if="showWhoisModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-4xl w-full mx-4 max-h-[90vh] overflow-y-auto">
          <div class="flex justify-between items-center p-6 border-b border-gray-200 dark:border-gray-700">
            <h3 class="text-xl font-bold text-gray-800 dark:text-white">WHOIS Information for {{ whoisInfo.domain }}</h3>
            <button
              @click="closeWhoisModal"
              class="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 text-2xl transition-colors"
            >
              ×
            </button>
          </div>

          <div v-if="whoisLoading" class="p-6 text-center">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
            <p class="mt-2 text-gray-600 dark:text-gray-400">Loading WHOIS information...</p>
          </div>

          <div v-else-if="whoisError" class="p-6">
            <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg">
              {{ whoisError }}
            </div>
          </div>

          <div v-else class="p-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 class="font-semibold text-gray-800 dark:text-white mb-3">Domain Information</h4>
                <div class="space-y-2">
                  <div>
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Domain:</span>
                    <span class="ml-2 font-mono text-gray-900 dark:text-white">{{ whoisInfo.domain }}</span>
                  </div>
                  <div v-if="whoisInfo.registrar">
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Registrar:</span>
                    <span class="ml-2 text-gray-900 dark:text-white">{{ whoisInfo.registrar }}</span>
                  </div>
                  <div v-if="whoisInfo.creation_date">
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Creation Date:</span>
                    <span class="ml-2 text-gray-900 dark:text-white">{{ whoisInfo.creation_date }}</span>
                  </div>
                  <div v-if="whoisInfo.expiration_date">
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Expiration Date:</span>
                    <span class="ml-2 text-gray-900 dark:text-white">{{ whoisInfo.expiration_date }}</span>
                  </div>
                  <div v-if="whoisInfo.updated_date">
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Updated Date:</span>
                    <span class="ml-2 text-gray-900 dark:text-white">{{ whoisInfo.updated_date }}</span>
                  </div>
                </div>
              </div>

              <div>
                <h4 class="font-semibold text-gray-800 dark:text-white mb-3">Technical Information</h4>
                <div class="space-y-2">
                  <div v-if="whoisInfo.status && whoisInfo.status.length > 0">
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Status:</span>
                    <div class="ml-2">
                      <span v-for="status in whoisInfo.status" :key="status"
                            class="inline-block bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-300 text-xs px-2 py-1 rounded mr-1 mb-1">
                        {{ status }}
                      </span>
                    </div>
                  </div>
                  <div v-if="whoisInfo.name_servers && whoisInfo.name_servers.length > 0">
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Name Servers:</span>
                    <div class="ml-2">
                      <div v-for="ns in whoisInfo.name_servers" :key="ns" class="font-mono text-sm text-gray-900 dark:text-white">
                        {{ ns }}
                      </div>
                    </div>
                  </div>
                  <div v-if="whoisInfo.admin_contact">
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Admin Contact:</span>
                    <span class="ml-2 text-gray-900 dark:text-white">{{ whoisInfo.admin_contact }}</span>
                  </div>
                  <div v-if="whoisInfo.tech_contact">
                    <span class="text-sm font-medium text-gray-600 dark:text-gray-400">Tech Contact:</span>
                    <span class="ml-2 text-gray-900 dark:text-white">{{ whoisInfo.tech_contact }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="whoisInfo.raw_data" class="mt-6">
              <h4 class="font-semibold text-gray-800 dark:text-white mb-3">Raw WHOIS Data</h4>
              <pre class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg text-xs overflow-x-auto text-gray-900 dark:text-white">{{ whoisInfo.raw_data }}</pre>
            </div>

            <div class="mt-6 text-sm text-gray-500 dark:text-gray-400">
              Last checked: {{ formatDate(whoisInfo.checked_at) }}
            </div>
          </div>

          <div class="flex justify-end p-6 border-t border-gray-200 dark:border-gray-700">
            <button
              @click="closeWhoisModal"
              class="px-4 py-2 bg-gray-600 dark:bg-gray-700 text-white rounded-lg hover:bg-gray-700 dark:hover:bg-gray-600 transition-colors"
            >
              Close
            </button>
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
      wsLoading: false,
      error: null,
      currentResult: null,
      bulkResult: null,
      wsProgress: null,
      domains: [],
      totalExtensions: 228,
      ws: null,
      showWhoisModal: false,
      whoisInfo: {},
      whoisLoading: false,
      whoisError: null,
      darkMode: false
    }
  },
  mounted() {
    this.loadDomains()
    this.initWebSocket()
    this.loadDarkMode()
  },
  beforeUnmount() {
    if (this.ws) {
      this.ws.close()
    }
  },
  methods: {
    loadDarkMode() {
      const saved = localStorage.getItem('darkMode')
      this.darkMode = saved ? JSON.parse(saved) : window.matchMedia('(prefers-color-scheme: dark)').matches
    },

    toggleDarkMode() {
      this.darkMode = !this.darkMode
      localStorage.setItem('darkMode', JSON.stringify(this.darkMode))
    },

    initWebSocket() {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsUrl = `${protocol}//${window.location.host}/ws`

      this.ws = new WebSocket(wsUrl)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
      }

      this.ws.onmessage = (event) => {
        const message = JSON.parse(event.data)
        this.handleWebSocketMessage(message)
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        this.error = 'WebSocket connection failed'
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
      }
    },

    handleWebSocketMessage(message) {
      switch (message.type) {
        case 'connected':
          console.log('WebSocket connected:', message.message)
          break
        case 'bulk_check_started':
          this.wsLoading = true
          this.error = null
          this.currentResult = null
          this.bulkResult = null
          this.wsProgress = null
          break
        case 'bulk_check_progress':
          this.wsProgress = message.data
          break
        case 'bulk_check_complete':
          this.wsProgress = message.data
          this.wsLoading = false
          this.domainNameInput = ''
          this.loadDomains()
          break
        case 'error':
          this.error = message.message
          this.wsLoading = false
          break
      }
    },

    checkAllExtensionsWebSocket() {
      if (!this.domainNameInput.trim() || !this.ws) return

      this.ws.send(JSON.stringify({
        type: 'check_all_extensions',
        data: {
          domain_name: this.domainNameInput.trim()
        }
      }))
    },

    async checkDomain() {
      if (!this.domainInput.trim()) return

      this.loading = true
      this.error = null
      this.currentResult = null
      this.bulkResult = null
      this.wsProgress = null

      try {
        const response = await axios.post('/api/check-domain', {
          domain: this.domainInput.trim()
        })

        if (response.data.success) {
          this.currentResult = response.data.data.domain
          this.domainInput = ''
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
      this.wsProgress = null

      try {
        const response = await axios.post('/api/check-all-extensions', {
          domain_name: this.domainNameInput.trim()
        })

        if (response.data.success) {
          this.bulkResult = response.data.data
          this.resultTab = 'available'
          this.domainNameInput = ''
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
          return 'text-green-600 dark:text-green-400'
        case 'Registered':
          return 'text-red-600 dark:text-red-400'
        case 'Error':
          return 'text-yellow-600 dark:text-yellow-400'
        default:
          return 'text-gray-600 dark:text-gray-400'
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
    },

    showWhoisInfo(domain) {
      this.showWhoisModal = true
      this.whoisInfo = {}
      this.whoisLoading = true
      this.whoisError = null

      axios.get(`/api/v1/domains/whois/${domain}`)
        .then(response => {
          if (response.data.success) {
            this.whoisInfo = response.data.data
          } else {
            this.whoisError = response.data.message || 'Failed to load WHOIS information'
          }
        })
        .catch(error => {
          this.whoisError = error.response?.data?.message || 'Network error occurred'
        })
        .finally(() => {
          this.whoisLoading = false
        })
    },

    closeWhoisModal() {
      this.showWhoisModal = false
    }
  }
}
</script>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}

/* Dark mode transitions */
* {
  transition: background-color 0.3s ease, border-color 0.3s ease, color 0.3s ease;
}

/* Custom scrollbar for dark mode */
.dark ::-webkit-scrollbar {
  width: 8px;
}

.dark ::-webkit-scrollbar-track {
  background: #374151;
}

.dark ::-webkit-scrollbar-thumb {
  background: #6B7280;
  border-radius: 4px;
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: #9CA3AF;
}
</style>
