package com.portfolio.blog.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.portfolio.blog.dto.AuthRequest;
import com.portfolio.blog.dto.AuthResponse;
import com.portfolio.blog.model.User;
import com.portfolio.blog.repository.UserRepository;
import com.portfolio.blog.security.JwtService;

import lombok.RequiredArgsConstructor;

@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

  private final AuthenticationManager authenticationManager;
  private final UserRepository userRepository;
  private final PasswordEncoder passwordEncoder;
  private final JwtService jwtService;
  private final UserDetailsService userDetailsService;

  @PostMapping("/register")
  public ResponseEntity<AuthResponse> register(@RequestBody AuthRequest request) {
    if (userRepository.findByUsername(request.getUsername()).isPresent()) {
      throw new RuntimeException("Username already exists");
    }

    User user = new User();
    user.setUsername(request.getUsername());
    user.setPassword(passwordEncoder.encode(request.getPassword()));
    user.setRole("ROLE_ADMIN");

    userRepository.save(user);

    UserDetails userDetails = userDetailsService.loadUserByUsername(request.getUsername());
    String token = jwtService.generateToken(userDetails);

    return ResponseEntity.ok(new AuthResponse(token));
  }

  @PostMapping("/login")
  public ResponseEntity<AuthResponse> login(@RequestBody AuthRequest request) {
    authenticationManager.authenticate(
        new UsernamePasswordAuthenticationToken(request.getUsername(), request.getPassword()));

    UserDetails userDetails = userDetailsService.loadUserByUsername(request.getUsername());
    String token = jwtService.generateToken(userDetails);

    return ResponseEntity.ok(new AuthResponse(token));
  }
}
