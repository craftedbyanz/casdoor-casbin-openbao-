// Casdoor Authentication Demo JavaScript

const API_BASE = 'http://localhost:8080';
let currentToken = null;

// Utility functions
function showResult(message, type = 'info') {
    const resultDiv = document.getElementById('result');
    resultDiv.innerHTML = `<div class="result ${type}">${message}</div>`;
}

function showError(message) {
    showResult(`❌ Error: ${message}`, 'error');
}

function showSuccess(message) {
    showResult(`✅ Success: ${message}`, 'success');
}

// Case 1: Direct Login
async function directLogin() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (!username || !password) {
        showError('Please enter username and password');
        return;
    }

    try {
        showResult('🔄 Logging in with username/password...');
        
        const response = await fetch(`${API_BASE}/api/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });

        const data = await response.json();

        if (response.ok) {
            localStorage.setItem('access_token', data.access_token);
            showSuccess('✅ Direct login successful! Redirecting to dashboard...');
            setTimeout(() => {
                window.location.href = '/dashboard.html';
            }, 1000);
        } else {
            showError(data.message || 'Login failed');
        }
    } catch (error) {
        showError(`Network error: ${error.message}`);
    }
}

// Case 2: OAuth/OIDC Login
async function oauthLogin() {
    try {
        showResult('🔄 Getting OAuth login URL...');
        
        const response = await fetch(`${API_BASE}/api/auth/oauth/login`);
        const data = await response.json();

        if (response.ok) {
            showResult(`🔗 Redirecting to OAuth login...`);
            localStorage.setItem('oauth_state', data.state);
            localStorage.setItem('auth_method', 'oauth');
            window.location.href = data.login_url;
        } else {
            showError(data.message || 'Failed to get OAuth URL');
        }
    } catch (error) {
        showError(`Network error: ${error.message}`);
    }
}

// Case 3: Microsoft SSO
async function microsoftSSO() {
    try {
        showResult('🔄 Getting Microsoft SSO URL...');
        
        const response = await fetch(`${API_BASE}/api/auth/microsoft/login`);
        const data = await response.json();

        if (response.ok) {
            showResult(`🔗 Redirecting to Microsoft login...`);
            localStorage.setItem('microsoft_state', data.state);
            localStorage.setItem('auth_method', 'microsoft');
            window.location.href = data.login_url;
        } else {
            showError(data.message || 'Failed to get Microsoft SSO URL');
        }
    } catch (error) {
        showError(`Network error: ${error.message}`);
    }
}

// Handle OAuth callback
async function handleCallback() {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');
    const state = urlParams.get('state');
    const authMethod = localStorage.getItem('auth_method') || 'oauth';

    if (code) {
        try {
            showResult(`🔄 Processing ${authMethod} callback...`);
            
            const response = await fetch(`${API_BASE}/api/auth/callback?code=${code}&state=${state}`);
            const data = await response.json();

            if (response.ok) {
                localStorage.setItem('access_token', data.access_token);
                showSuccess(`✅ ${authMethod.toUpperCase()} login successful! Redirecting to dashboard...`);
                // Clean localStorage
                localStorage.removeItem('oauth_state');
                localStorage.removeItem('microsoft_state');
                localStorage.removeItem('auth_method');
                setTimeout(() => {
                    window.location.href = '/dashboard.html';
                }, 1000);
            } else {
                showError(data.message || `${authMethod} callback failed`);
            }
        } catch (error) {
            showError(`Callback error: ${error.message}`);
        }
    }
}

// Show user information
async function showUserInfo() {
    if (!currentToken) return;

    try {
        const response = await fetch(`${API_BASE}/api/auth/me`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        const data = await response.json();

        if (response.ok) {
            document.getElementById('userInfo').style.display = 'block';
            document.getElementById('userDetails').innerHTML = `
                <p><strong>🆔 ID:</strong> ${data.id}</p>
                <p><strong>👤 Name:</strong> ${data.name}</p>
                <p><strong>📝 Display Name:</strong> ${data.display_name || 'N/A'}</p>
                <p><strong>📧 Email:</strong> ${data.email}</p>
                <p><strong>🔑 Is Admin:</strong> ${data.is_admin ? '✅ Yes' : '❌ No'}</p>
                <p><strong>🏢 Owner:</strong> ${data.owner || 'N/A'}</p>
                <details>
                    <summary><strong>🎫 Access Token</strong></summary>
                    <div class="token-display">${currentToken}</div>
                </details>
            `;
        } else {
            showError('Failed to get user info');
        }
    } catch (error) {
        showError(`Error getting user info: ${error.message}`);
    }
}

// Test protected endpoints
async function testProtected() {
    await testEndpoint('/api/protected', '🔒 Protected Resource');
}

async function testProfile() {
    await testEndpoint('/api/users/profile', '👤 User Profile');
}

async function testUsers() {
    await testEndpoint('/api/users', '👥 Users List (Admin Only)');
}

async function testSecrets() {
    await testEndpoint('/api/secrets', '🔐 Secrets (Cert Verification)');
}

// Transaction tests
async function testMyTransactions() {
    await testEndpoint('/api/transactions/my', '📊 My Transactions');
}

async function testAllTransactions() {
    await testEndpoint('/api/transactions', '📈 All Transactions (Admin Only)');
}

async function testCreateTransaction() {
    if (!currentToken) {
        showError('Please login first');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/api/transactions`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${currentToken}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                amount: 100.50,
                type: 'deposit',
                description: 'Test transaction from demo'
            })
        });

        const data = await response.json();
        const resultsDiv = document.getElementById('protectedResults');
        
        if (response.ok) {
            resultsDiv.innerHTML += `
                <div class="result success">
                    <strong>➕ Create Transaction:</strong> ✅ Success (${response.status})<br>
                    <details>
                        <summary>View Response</summary>
                        <pre>${JSON.stringify(data, null, 2)}</pre>
                    </details>
                </div>
            `;
        } else {
            resultsDiv.innerHTML += `
                <div class="result error">
                    <strong>➕ Create Transaction:</strong> ❌ ${response.status} - ${data.message || 'Failed'}
                </div>
            `;
        }
    } catch (error) {
        document.getElementById('protectedResults').innerHTML += `
            <div class="result error">
                <strong>➕ Create Transaction:</strong> ❌ Network error: ${error.message}
            </div>
        `;
    }
}

// Order tests
async function testMyOrders() {
    await testEndpoint('/api/orders/my', '📦 My Orders');
}

async function testAllOrders() {
    await testEndpoint('/api/orders', '📋 All Orders (Admin Only)');
}

async function testCreateOrder() {
    if (!currentToken) {
        showError('Please login first');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/api/orders`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${currentToken}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                product_name: 'Demo Product',
                quantity: 2,
                price: 25.99
            })
        });

        const data = await response.json();
        const resultsDiv = document.getElementById('protectedResults');
        
        if (response.ok) {
            resultsDiv.innerHTML += `
                <div class="result success">
                    <strong>🛍️ Create Order:</strong> ✅ Success (${response.status})<br>
                    <details>
                        <summary>View Response</summary>
                        <pre>${JSON.stringify(data, null, 2)}</pre>
                    </details>
                </div>
            `;
        } else {
            resultsDiv.innerHTML += `
                <div class="result error">
                    <strong>🛍️ Create Order:</strong> ❌ ${response.status} - ${data.message || 'Failed'}
                </div>
            `;
        }
    } catch (error) {
        document.getElementById('protectedResults').innerHTML += `
            <div class="result error">
                <strong>🛍️ Create Order:</strong> ❌ Network error: ${error.message}
            </div>
        `;
    }
}

async function testEndpoint(endpoint, name) {
    if (!currentToken) {
        showError('Please login first');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}${endpoint}`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        const data = await response.json();
        const resultsDiv = document.getElementById('protectedResults');
        
        if (response.ok) {
            resultsDiv.innerHTML += `
                <div class="result success">
                    <strong>${name}:</strong> ✅ Success (${response.status})<br>
                    <details>
                        <summary>View Response</summary>
                        <pre>${JSON.stringify(data, null, 2)}</pre>
                    </details>
                </div>
            `;
        } else {
            resultsDiv.innerHTML += `
                <div class="result error">
                    <strong>${name}:</strong> ❌ ${response.status} - ${data.message || 'Access denied'}
                </div>
            `;
        }
    } catch (error) {
        document.getElementById('protectedResults').innerHTML += `
            <div class="result error">
                <strong>${name}:</strong> ❌ Network error: ${error.message}
            </div>
        `;
    }
}

// Clear results
function clearResults() {
    document.getElementById('protectedResults').innerHTML = '';
}

// Logout
function logout() {
    currentToken = null;
    document.getElementById('userInfo').style.display = 'none';
    document.getElementById('result').innerHTML = '';
    document.getElementById('protectedResults').innerHTML = '';
    showSuccess('Logged out successfully');
}

// Initialize page
window.onload = function() {
    // Check if token exists in localStorage (from redirect)
    const storedToken = localStorage.getItem('access_token');
    if (storedToken && window.location.pathname === '/') {
        // Redirect to dashboard if token exists and we're on homepage
        window.location.href = '/dashboard.html';
        return;
    }
    
    // Handle OAuth callback if present
    if (window.location.search.includes('code=')) {
        handleCallback();
    }
};