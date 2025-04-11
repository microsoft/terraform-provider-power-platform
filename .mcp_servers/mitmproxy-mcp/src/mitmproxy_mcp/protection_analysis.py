import re
from typing import Any, Dict, List

# Known bot protection systems and their signatures
BOT_PROTECTION_SIGNATURES = {
    "Cloudflare": [
        r"cf-ray",  # Cloudflare Ray ID header
        r"__cf_bm",  # Cloudflare Bot Management cookie
        r"cf_clearance",  # Cloudflare challenge clearance cookie
        r"\"why_captcha\"",  # Common in Cloudflare challenge responses
        r"challenge-platform",  # Used in challenge scripts
        r"turnstile\.js",  # Cloudflare Turnstile
    ],
    "Akamai Bot Manager": [
        r"_abck=",  # Akamai Bot Manager cookie
        r"akam_", # Akamai cookie prefix
        r"bm_sz", # Bot Manager cookie
        r"sensor_data", # Bot detection data
    ],
    "PerimeterX": [
        r"_px\d?=",  # PerimeterX cookies
        r"px\.js", # PerimeterX script
        r"px-captcha", # PerimeterX captcha
    ],
    "DataDome": [
        r"datadome=",  # DataDome cookie
        r"datadome\.js", # DataDome script
        r"_dd_s",  # DataDome session cookie
    ],
    "reCAPTCHA": [
        r"google\.com/recaptcha",
        r"recaptcha\.net",
        r"g-recaptcha",
    ],
    "hCaptcha": [
        r"hcaptcha\.com",
        r"h-captcha",
    ],
    "Generic Bot Detection": [
        r"bot=",  # Generic bot cookie
        r"captcha", # Generic captcha reference
        r"challenge",  # Generic challenge term
        r"detected automated traffic", # Common message
        r"verify you are human", # Common message
    ]
}

def extract_javascript(html_content: str) -> List[Dict[str, Any]]:
    """
    Extract JavaScript from HTML content and provide basic analysis.
    Returns list of dictionaries with script info.
    """
    scripts = []
    
    # Extract inline scripts
    inline_pattern = r'<script[^>]*>(.*?)</script>'
    inline_scripts = re.findall(inline_pattern, html_content, re.DOTALL)
    
    for i, script in enumerate(inline_scripts):
        if len(script.strip()) > 0:
            script_info = {
                "type": "inline",
                "index": i,
                "size": len(script),
                "content": script if len(script) < 1000 else script[:1000] + "... [truncated]",
                "summary": analyze_script(script)
            }
            scripts.append(script_info)
    
    # Extract external script references
    src_pattern = r'<script[^>]*src=[\'"]([^\'"]+)[\'"][^>]*>'
    external_scripts = re.findall(src_pattern, html_content)
    
    for i, src in enumerate(external_scripts):
        script_info = {
            "type": "external",
            "index": i,
            "src": src,
            "suspicious": any(term in src.lower() for term in [
                "captcha", "challenge", "bot", "protect", "security", 
                "verify", "check", "shield", "defend", "guard"
            ])
        }
        scripts.append(script_info)
    
    return scripts

def analyze_script(script: str) -> Dict[str, Any]:
    """
    Analyze JavaScript content for common protection patterns.
    """
    analysis = {
        "potential_protection": False,
        "fingerprinting_indicators": [],
        "token_generation_indicators": [],
        "obfuscation_level": "none",
        "key_functions": []
    }
    
    # Check for fingerprinting techniques
    fingerprinting_patterns = [
        (r'navigator\.', "Browser navigator object"),
        (r'screen\.', "Screen properties"),
        (r'canvas', "Canvas fingerprinting"),
        (r'webgl', "WebGL fingerprinting"),
        (r'font', "Font enumeration"),
        (r'audio', "Audio fingerprinting"),
        (r'plugins', "Plugin enumeration"),
        (r'User-Agent', "User-Agent checking"),
        (r'platform', "Platform detection")
    ]
    
    for pattern, description in fingerprinting_patterns:
        if re.search(pattern, script, re.IGNORECASE):
            analysis["fingerprinting_indicators"].append(description)
    
    # Check for token generation
    token_patterns = [
        (r'(token|captcha|challenge|clearance)', "Token/challenge reference"),
        (r'(generate|calculate|compute)', "Computation terms"),
        (r'(Math\.random|crypto)', "Random generation"),
        (r'(cookie|setCookie|document\.cookie)', "Cookie manipulation"),
        (r'(xhr|XMLHttpRequest|fetch)', "Request sending")
    ]
    
    for pattern, description in token_patterns:
        if re.search(pattern, script, re.IGNORECASE):
            analysis["token_generation_indicators"].append(description)
    
    # Check for common obfuscation techniques
    if len(re.findall(r'eval\(', script)) > 3:
        analysis["obfuscation_level"] = "high"
    elif len(re.findall(r'\\x[0-9a-f]{2}', script)) > 10:
        analysis["obfuscation_level"] = "high"
    elif len(re.findall(r'String\.fromCharCode', script)) > 3:
        analysis["obfuscation_level"] = "high"
    elif re.search(r'function\(\w{1,2},\w{1,2},\w{1,2}\)\{', script):
        analysis["obfuscation_level"] = "medium"
    elif sum(1 for c in script if c == ';') > len(script) / 10:
        analysis["obfuscation_level"] = "medium"
    elif sum(len(w) > 30 for w in re.findall(r'\w+', script)) > 10:
        analysis["obfuscation_level"] = "medium"
    
    # Extract potential key function names
    function_pattern = r'function\s+(\w+)\s*\('
    functions = re.findall(function_pattern, script)
    
    suspicious_terms = ["challenge", "token", "captcha", "verify", "bot", "check", "security"]
    for func in functions:
        if any(term in func.lower() for term in suspicious_terms):
            analysis["key_functions"].append(func)
    
    # Determine if this is potentially protection-related
    analysis["potential_protection"] = (
        len(analysis["fingerprinting_indicators"]) > 2 or
        len(analysis["token_generation_indicators"]) > 2 or
        analysis["obfuscation_level"] != "none" or
        len(analysis["key_functions"]) > 0
    )
    
    return analysis

def analyze_cookies(headers: Dict[str, str]) -> List[Dict[str, Any]]:
    """
    Analyze cookies for common protection-related patterns.
    """
    cookie_header = headers.get("Cookie", "") or headers.get("Set-Cookie", "")
    if not cookie_header:
        return []
    
    # Split multiple cookies
    cookies = []
    for cookie_str in cookie_header.split(";"):
        parts = cookie_str.strip().split("=", 1)
        if len(parts) == 2:
            name, value = parts
            cookie = {
                "name": name.strip(),
                "value": value.strip() if len(value.strip()) < 50 else value.strip()[:50] + "... [truncated]",
                "protection_related": False,
                "vendor": "unknown"
            }
            
            # Check if this is a known protection cookie
            for vendor, signatures in BOT_PROTECTION_SIGNATURES.items():
                for sig in signatures:
                    if re.search(sig, name, re.IGNORECASE):
                        cookie["protection_related"] = True
                        cookie["vendor"] = vendor
                        break
                if cookie["protection_related"]:
                    break
            
            cookies.append(cookie)
    
    return cookies

def identify_protection_system(flow) -> List[Dict[str, Any]]:
    """
    Identify potential bot protection systems based on signatures.
    """
    protections = []
    
    # Combine all searchable content
    searchable_content = ""
    # Add request headers
    for k, v in flow.request.headers.items():
        searchable_content += f"{k}: {v}\n"
    
    # Check response if available
    if flow.response:
        # Add response headers
        for k, v in flow.response.headers.items():
            searchable_content += f"{k}: {v}\n"
        
        # Add response content if it's text
        content_type = flow.response.headers.get("Content-Type", "")
        if "text" in content_type or "javascript" in content_type or "json" in content_type:
            try:
                searchable_content += flow.response.content.decode('utf-8', errors='ignore')
            except Exception:
                pass
    
    # Check for protection signatures
    for vendor, signatures in BOT_PROTECTION_SIGNATURES.items():
        matches = []
        for sig in signatures:
            if re.search(sig, searchable_content, re.IGNORECASE):
                matches.append(sig)
        
        if matches:
            protections.append({
                "vendor": vendor,
                "confidence": len(matches) / len(signatures) * 100,
                "matching_signatures": matches
            })
    
    return sorted(protections, key=lambda x: x["confidence"], reverse=True)

def analyze_response_for_challenge(flow) -> Dict[str, Any]:
    """
    Analyze a response to determine if it contains a challenge.
    """
    if not flow.response:
        return {"is_challenge": False}
    
    result = {
        "is_challenge": False,
        "challenge_indicators": [],
        "status_code": flow.response.status_code,
        "challenge_type": "unknown"
    }
    
    # Check status code
    if flow.response.status_code in [403, 429, 503]:
        result["challenge_indicators"].append(f"Suspicious status code: {flow.response.status_code}")
    
    # Check for challenge headers
    challenge_headers = {
        "cf-mitigated": "Cloudflare mitigation",
        "cf-chl-bypass": "Cloudflare challenge bypass",
        "x-datadome": "DataDome protection",
        "x-px": "PerimeterX",
        "x-amz-captcha": "AWS WAF Captcha"
    }
    
    for header, description in challenge_headers.items():
        if any(h.lower() == header.lower() for h in flow.response.headers.keys()):
            result["challenge_indicators"].append(f"Challenge header: {description}")
    
    # Check for challenge content patterns
    content = flow.response.content.decode('utf-8', errors='ignore')
    challenge_patterns = [
        (r'captcha', "CAPTCHA"),
        (r'challenge', "Challenge term"),
        (r'blocked', "Blocking message"),
        (r'verify.*human', "Human verification"),
        (r'suspicious.*activity', "Suspicious activity message"),
        (r'security.*check', "Security check message"),
        (r'ddos', "DDoS protection message"),
        (r'automated.*request', "Automated request detection")
    ]
    
    for pattern, description in challenge_patterns:
        if re.search(pattern, content, re.IGNORECASE):
            result["challenge_indicators"].append(f"Content indicator: {description}")
    
    # Determine if this is a challenge response
    result["is_challenge"] = len(result["challenge_indicators"]) > 0
    
    # Determine challenge type
    if "CAPTCHA" in " ".join(result["challenge_indicators"]):
        result["challenge_type"] = "captcha"
    elif "JavaScript" in content and result["is_challenge"]:
        result["challenge_type"] = "javascript"
    elif result["is_challenge"]:
        result["challenge_type"] = "other"
    
    return result

def generate_suggestions(analysis: Dict[str, Any]) -> List[str]:
    """
    Generate remediation suggestions based on the protection analysis.
    """
    suggestions = []
    
    # Check if any protection system was detected
    if analysis.get("protection_systems"):
        top_system = analysis["protection_systems"][0]["vendor"]
        confidence = analysis["protection_systems"][0]["confidence"]
        
        if confidence > 50:
            suggestions.append(f"Detected {top_system} with {confidence:.1f}% confidence.")
            
            # Add vendor-specific suggestions
            if "Cloudflare" in top_system:
                suggestions.append("Cloudflare often uses JavaScript challenges. Check for cf_clearance cookie.")
                suggestions.append("Consider using proven techniques like cfscrape or cloudscraper libraries.")
            elif "Akamai" in top_system:
                suggestions.append("Akamai uses sensor_data for browser fingerprinting.")
                suggestions.append("Focus on _abck cookie which contains browser verification data.")
            elif "PerimeterX" in top_system:
                suggestions.append("PerimeterX relies on JavaScript execution and browser fingerprinting.")
                suggestions.append("Look for _px cookies which are essential for session validation.")
            elif "DataDome" in top_system:
                suggestions.append("DataDome uses advanced behavioral and fingerprinting techniques.")
                suggestions.append("The datadome cookie is critical for maintaining sessions.")
            elif "CAPTCHA" in top_system:
                suggestions.append("This site uses CAPTCHA challenges which may require manual solving or specialized services.")
    
    # Add suggestions based on challenge type
    if analysis.get("challenge_analysis", {}).get("is_challenge", False):
        challenge_type = analysis.get("challenge_analysis", {}).get("challenge_type", "unknown")
        
        if challenge_type == "javascript":
            suggestions.append("This response contains a JavaScript challenge that must be solved.")
            suggestions.append("Consider using a headless browser to execute the challenge JavaScript.")
            
            # If we have script analysis, add more specific suggestions
            if "scripts" in analysis:
                obfuscated_scripts = [s for s in analysis["scripts"] if s.get("summary", {}).get("obfuscation_level") in ["medium", "high"]]
                if obfuscated_scripts:
                    suggestions.append(f"Found {len(obfuscated_scripts)} obfuscated script(s) that likely contain challenge logic.")
                
                fingerprinting_scripts = [s for s in analysis["scripts"] if s.get("summary", {}).get("fingerprinting_indicators")]
                if fingerprinting_scripts:
                    techniques = set()
                    for script in fingerprinting_scripts:
                        techniques.update(script.get("summary", {}).get("fingerprinting_indicators", []))
                    suggestions.append(f"Detected browser fingerprinting techniques: {', '.join(techniques)}.")
                    
        elif challenge_type == "captcha":
            suggestions.append("This response contains a CAPTCHA challenge.")
            suggestions.append("Consider using a CAPTCHA solving service or manual intervention.")
    
    # Check for important cookies
    protection_cookies = [c for c in analysis.get("response_cookies", []) if c.get("protection_related")]
    if protection_cookies:
        cookie_names = [c["name"] for c in protection_cookies]
        suggestions.append(f"Important protection cookies to maintain: {', '.join(cookie_names)}.")
    
    # General suggestions
    if analysis.get("protection_systems") or analysis.get("challenge_analysis", {}).get("is_challenge", False):
        suggestions.append("General recommendations:")
        suggestions.append("- Maintain consistent User-Agent between requests")
        suggestions.append("- Preserve all cookies from the session")
        suggestions.append("- Add appropriate referer and origin headers")
        suggestions.append("- Consider adding delays between requests to avoid rate limiting")
        suggestions.append("- Use rotating IP addresses if available")
    
    return suggestions
