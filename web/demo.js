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
            currentToken = data.access_token;
            showSuccess('✅ Direct login successful!');
            await showUserInfo();
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
                currentToken = data.access_token;
                showSuccess(`✅ ${authMethod.toUpperCase()} login successful!`);
                await showUserInfo();
                // Clean URL and localStorage
                window.history.replaceState({}, document.title, window.location.pathname);
                localStorage.removeItem('oauth_state');
                localStorage.removeItem('microsoft_state');
                localStorage.removeItem('auth_method');
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
    // Handle OAuth callback if present
    if (window.location.search.includes('code=')) {
        handleCallback();
    }
};