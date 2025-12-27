package com.portfolio.blog.service;

import java.util.Locale;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.portfolio.blog.dto.BlogRequest;
import com.portfolio.blog.model.BlogPost;
import com.portfolio.blog.repository.BlogRepository;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class BlogService {

  private final BlogRepository blogRepository;

  public Page<BlogPost> getAllPosts(Pageable pageable) {
    return blogRepository.findAll(pageable);
  }

  public Page<BlogPost> getPublishedPosts(Pageable pageable) {
    return blogRepository.findByPublishedTrue(pageable);
  }

  public Page<BlogPost> searchPosts(String query, Pageable pageable) {
    return blogRepository.findByTitleContainingIgnoreCaseOrTagsContainingIgnoreCase(query, query, pageable);
  }

  public BlogPost getPostBySlug(String slug) {
    return blogRepository.findBySlug(slug)
        .orElseThrow(() -> new RuntimeException("Blog post not found with slug: " + slug));
  }

  @Transactional
  public BlogPost createPost(BlogRequest request) {
    BlogPost post = new BlogPost();
    post.setTitle(request.getTitle());
    post.setContent(request.getContent());
    post.setCoverImage(request.getCoverImage());
    post.setTags(request.getTags());
    post.setPublished(request.isPublished());

    String slug = generateSlug(request.getTitle());
    if (blogRepository.existsBySlug(slug)) {
      slug = slug + "-" + System.currentTimeMillis();
    }
    post.setSlug(slug);

    return blogRepository.save(post);
  }

  @Transactional
  public BlogPost updatePost(Long id, BlogRequest request) {
    BlogPost post = blogRepository.findById(id)
        .orElseThrow(() -> new RuntimeException("Blog post not found with id: " + id));

    post.setTitle(request.getTitle());
    post.setContent(request.getContent());
    post.setCoverImage(request.getCoverImage());
    post.setTags(request.getTags());
    post.setPublished(request.isPublished());

    // Optionally update slug if title changes, but often it's better to keep slug
    // stable for SEO.
    // For now, we will NOT update the slug automatically to preserve links.

    return blogRepository.save(post);
  }

  @Transactional
  public void deletePost(Long id) {
    if (!blogRepository.existsById(id)) {
      throw new RuntimeException("Blog post not found with id: " + id);
    }
    blogRepository.deleteById(id);
  }

  private String generateSlug(String title) {
    return title.toLowerCase(Locale.ROOT)
        .replaceAll("[^a-z0-9\\s-]", "")
        .replaceAll("\\s+", "-");
  }
}
