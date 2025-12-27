package com.portfolio.blog.controller;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.portfolio.blog.dto.BlogRequest;
import com.portfolio.blog.model.BlogPost;
import com.portfolio.blog.service.BlogService;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;

@RestController
@RequestMapping("/api/blogs")
@RequiredArgsConstructor
@CrossOrigin(origins = "*") // Configure properly for production later
public class BlogController {

  private final BlogService blogService;

  @GetMapping
  public ResponseEntity<Page<BlogPost>> getAllPosts(
      @PageableDefault(size = 10, sort = "createdAt", direction = Sort.Direction.DESC) Pageable pageable,
      @RequestParam(required = false) String search,
      @RequestParam(defaultValue = "false") boolean publicOnly) {

    if (search != null && !search.isEmpty()) {
      return ResponseEntity.ok(blogService.searchPosts(search, pageable));
    }

    if (publicOnly) {
      return ResponseEntity.ok(blogService.getPublishedPosts(pageable));
    }

    return ResponseEntity.ok(blogService.getAllPosts(pageable));
  }

  @GetMapping("/{slug}")
  public ResponseEntity<BlogPost> getPostBySlug(@PathVariable String slug) {
    return ResponseEntity.ok(blogService.getPostBySlug(slug));
  }

  @PostMapping
  public ResponseEntity<BlogPost> createPost(@Valid @RequestBody BlogRequest request) {
    return ResponseEntity.ok(blogService.createPost(request));
  }

  @PutMapping("/{id}")
  public ResponseEntity<BlogPost> updatePost(@PathVariable Long id, @Valid @RequestBody BlogRequest request) {
    return ResponseEntity.ok(blogService.updatePost(id, request));
  }

  @DeleteMapping("/{id}")
  public ResponseEntity<Void> deletePost(@PathVariable Long id) {
    blogService.deletePost(id);
    return ResponseEntity.noContent().build();
  }
}
